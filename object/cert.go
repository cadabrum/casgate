// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/casdoor/casdoor/orm"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/core"
)

const (
	scopeCertJWT    = "JWT"
	scopeCertCACert = "CA Certificate"
	scopeClientCert = "Client Certificate"
)

type Cert struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName     string `xorm:"varchar(100)" json:"displayName"`
	Scope           string `xorm:"varchar(100)" json:"scope"`
	Type            string `xorm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `xorm:"varchar(100)" json:"cryptoAlgorithm"`
	BitSize         int    `json:"bitSize"`
	ExpireInYears   int    `json:"expireInYears"`

	Certificate string `xorm:"mediumtext" json:"certificate"`
	PrivateKey  string `xorm:"mediumtext" json:"privateKey"`
}

var ErrCertDoesNotExist = errors.New(fmt.Sprintf("certificate does not exist"))
var ErrCertInvalidScope = errors.New(fmt.Sprintf("invalid certificate scope"))

func GetMaskedCert(cert *Cert) *Cert {
	if cert == nil {
		return nil
	}

	if cert.PrivateKey != "" {
		cert.PrivateKey = "***"
	}

	return cert
}

func GetMaskedCerts(certs []*Cert, err error) ([]*Cert, error) {
	if err != nil {
		return nil, err
	}

	for _, cert := range certs {
		cert = GetMaskedCert(cert)
	}
	return certs, nil
}

func GetCertCount(owner, field, value string) (int64, error) {
	session := orm.GetSession("", -1, -1, field, value, "", "")
	return session.Where("owner = ? or owner = ? ", "admin", owner).Count(&Cert{})
}

func GetCerts(owner string) ([]*Cert, error) {
	certs := []*Cert{}
	err := orm.AppOrmer.Engine.Where("owner = ? or owner = ? ", "admin", owner).Desc("created_time").Find(&certs, &Cert{})
	if err != nil {
		return certs, err
	}

	return certs, nil
}

func GetPaginationCerts(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*Cert, error) {
	certs := []*Cert{}
	session := orm.GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Where("owner = ? or owner = ? ", "admin", owner).Find(&certs)
	if err != nil {
		return certs, err
	}

	return certs, nil
}

func GetGlobalCertsCount(field, value string) (int64, error) {
	session := orm.GetSession("", -1, -1, field, value, "", "")
	return session.Count(&Cert{})
}

func GetGlobleCerts(owner string) ([]*Cert, error) {
	certs := []*Cert{}
	session := orm.GetSession("", -1, -1, "", "", "", "")

	var err error
	if owner != "" {
		err = session.Where("owner = ? or owner = ? ", "admin", owner).Desc("created_time").Find(&certs)
	} else {
		session.Desc("created_time").Find(&certs)
	}

	if err != nil {
		return certs, err
	}

	return certs, nil
}

func GetPaginationGlobalCerts(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*Cert, error) {
	certs := []*Cert{}
	session := orm.GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&certs)
	if err != nil {
		return certs, err
	}

	return certs, nil
}

func getCert(owner string, name string) (*Cert, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	cert := Cert{Owner: owner, Name: name}
	existed, err := orm.AppOrmer.Engine.Get(&cert)
	if err != nil {
		return &cert, err
	}

	if existed {
		return &cert, nil
	} else {
		return nil, nil
	}
}

func getCertByName(name string) (*Cert, error) {
	if name == "" {
		return nil, nil
	}

	cert := Cert{Name: name}
	existed, err := orm.AppOrmer.Engine.Get(&cert)
	if err != nil {
		return &cert, nil
	}

	if existed {
		return &cert, nil
	} else {
		return nil, nil
	}
}

func GetTlsConfigForCert(name string) (*tls.Config, error) {
	cert, err := getCertByName(name)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, ErrCertDoesNotExist
	}

	ca := x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM([]byte(cert.Certificate)); !ok {
		return nil, ErrX509CertsPEMParse
	}

	return &tls.Config{RootCAs: ca}, nil
}

func GetCert(id string) (*Cert, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getCert(owner, name)
}

func UpdateCert(id string, cert *Cert) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if c, err := getCert(owner, name); err != nil {
		return false, err
	} else if c == nil {
		return false, nil
	}

	if name != cert.Name {
		err := certChangeTrigger(name, cert.Name)
		if err != nil {
			return false, err
		}
	}
	session := orm.AppOrmer.Engine.ID(core.PK{owner, name}).AllCols()
	if cert.PrivateKey == "***" {
		session.Omit("private_key")
	}
	affected, err := session.Update(cert)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddCert(cert *Cert) (bool, error) {
	if cert.Scope == scopeCertJWT && (cert.Certificate == "" || cert.PrivateKey == "") {
		certificate, privateKey := generateRsaKeys(cert.BitSize, cert.ExpireInYears, cert.Name, cert.Owner)
		cert.Certificate = certificate
		cert.PrivateKey = privateKey
	}

	affected, err := orm.AppOrmer.Engine.Insert(cert)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteCert(cert *Cert) (bool, error) {
	affected, err := orm.AppOrmer.Engine.ID(core.PK{cert.Owner, cert.Name}).Delete(&Cert{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (p *Cert) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}

func getCertByApplication(application *Application) (*Cert, error) {
	if application.Cert != "" {
		return getCertByName(application.Cert)
	} else {
		return GetDefaultCert()
	}
}

func GetDefaultCert() (*Cert, error) {
	return getCert("admin", "cert-built-in")
}

func certChangeTrigger(oldName string, newName string) error {
	session := orm.AppOrmer.Engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	application := new(Application)
	application.Cert = newName
	_, err = session.Where("cert=?", oldName).Update(application)
	if err != nil {
		return err
	}

	return session.Commit()
}
