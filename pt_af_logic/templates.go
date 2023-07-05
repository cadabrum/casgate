package pt_af_logic

const partnerCreateAccountsSubjTmpl = `[PT LMP] Accounts for {{ .ClientName }}`

const partnerCreateAccountsBodyTmpl = `
<p>В PT AF создано изолированное пространство для клиента <a href="{{ .ClientURL }}">{{ .ClientDisplayName }}</a>.</p>

<p>Строка подключения для агента: {{ .ConnectionString }}</p>

<p>Ссылка для входа: {{ .PTAFLoginLink }}<br>
При первом входе потребуется сменить пароль.</p>

<p>Сервисная учётная запись<br>
Логин: {{ .ServiceUserName }}<br>
Временный пароль: {{ .ServiceUserPwd }}</p>

<p>Пользовательская учётная запись<br>
Логин: {{ .UserROName }}<br>
Временный пароль: {{ .UserROPwd }}</p>
`

const partnerCreatedSubjTmpl = `[PT LMP] Partner {{ .PartnerName }} registered`

const partnerCreatedBodyTmpl = `
<p>Партнёр <a href="{{ .PartnerURL }}">{{ .PartnerDisplayName }}</a> зарегистрировался на портале.<br>
Для подтверждения регистрации включите опцию "Is admin" и отключите "Is forbidden" в аккаунте партнёра: <a href="{{ .PartnerAccount }}">{{ .PartnerUserName }}</a>.</p>
`

const partnerConfirmedBodyTmpl = `
<p>Ваша регистрация на портале подтверждена.</p>
<p>Логин: {{ .PartnerUserName }}<br>
Ссылка для входа на портал: <a href="{{ .PartnerLoginURL }}">{{ .PartnerLoginURL }}</a></p>
`

const SubscriptionUpdatedSubjTmpl = `[PT LMP] Обновлена подписка {{ .SubscriptionName }} (партнёр {{ .PartnerName }}, заказчик {{ .ClientName }}){{if ne .OldSubscriptionStatus .SubscriptionStatus}}: {{ .SubscriptionStatus }}{{end}}`

const SubscriptionUpdatedBodyTmpl = `
<p>Партнёр: <a href="{{ .PartnerURL }}">{{ .PartnerDisplayName }}</a><br>
Идентификатор: <a href="{{ .SubscriptionURL }}">{{ .SubscriptionName }}</a><br>
Заказчик:<a href="{{ .ClientURL }}">{{ .ClientDisplayName }}</a><br>
Тарифный план: <a href="{{ .OldPlanURL }}">{{ .OldPlanDisplayName }}</a>{{if ne .OldPlanDisplayName .PlanDisplayName}} -> <a href="{{ .PlanURL }}">{{ .PlanDisplayName }}</a>{{end}}<br>
Скидка: {{ .OldSubscriptionDiscount }}{{if ne .OldSubscriptionDiscount .SubscriptionDiscount}} -> {{ .SubscriptionDiscount }}{{end}}<br>
Дата начала: {{ .OldSubscriptionStartDate }}{{if ne .OldSubscriptionStartDate .SubscriptionStartDate}} -> {{ .SubscriptionStartDate }}{{end}}<br>
Дата окончания: {{ .OldSubscriptionEndDate }}{{if ne .OldSubscriptionEndDate .SubscriptionEndDate}} -> {{ .SubscriptionEndDate }}{{end}}<br>
Статус: {{ .OldSubscriptionStatus }}{{if ne .OldSubscriptionStatus .SubscriptionStatus}} -> {{ .SubscriptionStatus }}{{end}}<br>
Описание: {{ .OldSubscriptionDescription }}{{if ne .OldSubscriptionDescription .SubscriptionDescription }} -> {{ .SubscriptionDescription }}{{end}}<br>
Комментарий: {{ .OldSubscriptionComment }}{{if ne .OldSubscriptionComment .SubscriptionComment }} -> {{ .SubscriptionComment }}{{end}}<br>
Создатель: <a href="{{ .SubscriptionCreatorURL }}">{{ .SubscriptionCreator }}</a><br>
Последняя смена статуса: <a href="{{ .SubscriptionMoverURL }}">{{ .SubscriptionMover }}</a><br>
Время последней смены статуса: {{ .SubscriptionMoveTime }}<br>
{{ if eq .OldSubscriptionStatus .SubscriptionStatus }}
Последнее изменение: <a href="{{ .SubscriptionEditorURL }}">{{ .SubscriptionEditor }}</a><br>
Время последнего изменения: {{ .SubscriptionEditTime }}{{end}}</p>
`

const SubscriptionUpdatedPartnerSubjTmpl = `[PT LMP] Обновлена подписка {{ .SubscriptionName }} (заказчик {{ .ClientName }}): {{ .SubscriptionStatus }}`

const SubscriptionUpdatedPartnerBodyTmpl = `
<p>Идентификатор: <a href="{{ .SubscriptionURL }}">{{ .SubscriptionName }}</a><br>
Заказчик:<a href="{{ .ClientURL }}">{{ .ClientDisplayName }}</a><br>
Тарифный план: <a href="{{ .OldPlanURL }}">{{ .OldPlanDisplayName }}</a>{{if ne .OldPlanDisplayName .PlanDisplayName}} -> <a href="{{ .PlanURL }}">{{ .PlanDisplayName }}</a>{{end}}<br>
Скидка: {{ .OldSubscriptionDiscount }}{{if ne .OldSubscriptionDiscount .SubscriptionDiscount}} -> {{ .SubscriptionDiscount }}{{end}}<br>
Дата начала: {{ .OldSubscriptionStartDate }}{{if ne .OldSubscriptionStartDate .SubscriptionStartDate}} -> {{ .SubscriptionStartDate }}{{end}}<br>
Дата окончания: {{ .OldSubscriptionEndDate }}{{if ne .OldSubscriptionEndDate .SubscriptionEndDate}} -> {{ .SubscriptionEndDate }}{{end}}<br>
Статус: {{ .OldSubscriptionStatus }}{{if ne .OldSubscriptionStatus .SubscriptionStatus}} -> {{ .SubscriptionStatus }}{{end}}<br>
Описание: {{ .OldSubscriptionDescription }}{{if ne .OldSubscriptionDescription .SubscriptionDescription }} -> {{ .SubscriptionDescription }}{{end}}<br>
Комментарий: {{ .OldSubscriptionComment }}{{if ne .OldSubscriptionComment .SubscriptionComment }} -> {{ .SubscriptionComment }}{{end}}<br>
Создатель: <a href="{{ .SubscriptionCreatorURL }}">{{ .SubscriptionCreator }}</a><br>
Время последней смены статуса: {{ .SubscriptionMoveTime }}</p>
`