package amocrm

const (
	AccountNotFoundCode          = 101
	BodyMustBeJSONCode           = 102
	InvalidRequestParametersCode = 103
	InvalidRequestMethodCode     = 104

	InvalidCredentialsCode    = 110
	CaptchaInputRequiredCode  = 111
	DisabledAccountCode       = 112
	UnauthorizedIpAddressCode = 113

	ContactAddEmptyArrayCode                = 201
	ContactAddInsufficientAccessRightsCode  = 202
	ContactAddCustomFieldsInternalErrorCode = 203
	ContactAddCustomFieldNotFoundCode       = 204
	ContactAddNotProcessedCode              = 205

	ContactsEmptyRequestCode         = 206
	ContactsInvalidRequestMethodCode = 207

	ContactUpdateEmptyArrayCode                 = 208
	ContactUpdatedRequiredParametersMissingCode = 209
	ContactUpdateCustomFieldsInternalErrorCode  = 210
	ContactUpdateCustomFieldNotFoundCode        = 211
	ContactUpdateNotProcessedCode               = 212

	LeadAddEmptyArrayCode                   = 213
	LeadEmptyRequestCode                    = 214
	LeadInvalidRequestMethodCode            = 215
	LeadUpdateEmptyArrayCode                = 216
	LeadUpdateRequiredParametersMissingCode = 217

	ContactGetSearchErrorCode = 219

	LeadCustomFieldInvalidIDCode = 240

	TooManyLinkedEntitiesCode = 330

	InvalidRequestCode          = 400
	AccountNotFoundOnServerCode = 401
	SubscriptionExpireCode      = 402
	AccountBlockedCode          = 403

	RateLimitExceededCode = 429

	NoContentCode = 2002

	AccountNotFound          = "account_not_found"
	BodyMustBeJSON           = "body_must_be_json"
	InvalidRequestParameters = "invalid_request_params"
	InvalidRequestMethod     = "invalid_request_method"

	InvalidCredentials    = "invalid_credentials"
	CaptchaInputRequired  = "captcha_input_required"
	DisabledAccount       = "disabled_account"
	UnauthorizedIpAddress = "unauthorized_ip_address"

	ContactAddEmptyArray                = "contact_add_empty_array"
	ContactAddInsufficientAccessRights  = "contact_add_insufficient_access_rights"
	ContactAddCustomFieldsInternalError = "contact_add_custom_fields_internal_error"
	ContactAddCustomFieldNotFound       = "contact_add_custom_field_not_found"
	ContactAddNotProcessed              = "contact_add_not_processed"

	ContactsEmptyRequest         = "contacts_empty_request"
	ContactsInvalidRequestMethod = "contacts_invalid_request_method"

	ContactUpdateEmptyArray                 = "contact_update_empty_array"
	ContactUpdatedRequiredParametersMissing = "contact_update_required_params_missing"
	ContactUpdateCustomFieldsInternalError  = "contact_update_custom_fields_internal_error"
	ContactUpdateCustomFieldNotFound        = "contact_update_custom_field_not_found"
	ContactUpdateNotProcessed               = "contact_update_not_processed"

	LeadAddEmptyArray                   = "lead_add_empty_array"
	LeadEmptyRequest                    = "leads_empty_request"
	LeadInvalidRequestMethod            = "leads_invalid_request_method"
	LeadUpdateEmptyArray                = "lead_update_empty_array"
	LeadUpdateRequiredParametersMissing = "lead_update_required_params_missing"

	ContactGetSearchError = "contact_get_search_error"

	LeadCustomFieldInvalidID = "leads_custom_field_invalid_id"

	TooManyLinkedEntities = "too_many_linked_entities"

	InvalidRequest          = "invalid_request"
	AccountNotFoundOnServer = "account_not_found_on_server"
	SubscriptionExpire      = "subscription_expired"
	AccountBlocked          = "account_blocked"

	RateLimitExceeded = "rate_limit_exceeded"

	NoContent = "no_content"
)

var (
	AmoErrorTypeMap = map[int]string{
		AccountNotFoundCode:          AccountNotFound,
		BodyMustBeJSONCode:           BodyMustBeJSON,
		InvalidRequestParametersCode: InvalidRequestParameters,
		InvalidRequestMethodCode:     InvalidRequestMethod,

		InvalidCredentialsCode:    InvalidCredentials,
		CaptchaInputRequiredCode:  CaptchaInputRequired,
		DisabledAccountCode:       DisabledAccount,
		UnauthorizedIpAddressCode: UnauthorizedIpAddress,

		ContactAddEmptyArrayCode:                ContactAddEmptyArray,
		ContactAddInsufficientAccessRightsCode:  ContactAddInsufficientAccessRights,
		ContactAddCustomFieldsInternalErrorCode: ContactAddCustomFieldsInternalError,
		ContactAddCustomFieldNotFoundCode:       ContactAddCustomFieldNotFound,
		ContactAddNotProcessedCode:              ContactAddNotProcessed,

		ContactsEmptyRequestCode:         ContactsEmptyRequest,
		ContactsInvalidRequestMethodCode: ContactsInvalidRequestMethod,

		ContactUpdateEmptyArrayCode:                 ContactUpdateEmptyArray,
		ContactUpdatedRequiredParametersMissingCode: ContactUpdatedRequiredParametersMissing,
		ContactUpdateCustomFieldsInternalErrorCode:  ContactUpdateCustomFieldsInternalError,
		ContactUpdateCustomFieldNotFoundCode:        ContactUpdateCustomFieldNotFound,
		ContactUpdateNotProcessedCode:               ContactUpdateNotProcessed,

		LeadAddEmptyArrayCode:                   LeadAddEmptyArray,
		LeadEmptyRequestCode:                    LeadEmptyRequest,
		LeadInvalidRequestMethodCode:            LeadInvalidRequestMethod,
		LeadUpdateEmptyArrayCode:                LeadUpdateEmptyArray,
		LeadUpdateRequiredParametersMissingCode: LeadUpdateRequiredParametersMissing,

		ContactGetSearchErrorCode: ContactGetSearchError,

		LeadCustomFieldInvalidIDCode: LeadCustomFieldInvalidID,

		TooManyLinkedEntitiesCode: TooManyLinkedEntities,

		InvalidRequestCode:          InvalidRequest,
		AccountNotFoundOnServerCode: AccountNotFoundOnServer,
		SubscriptionExpireCode:      SubscriptionExpire,
		AccountBlockedCode:          AccountBlocked,

		RateLimitExceededCode: RateLimitExceeded,

		NoContentCode: NoContent,
	}
)

type (
	AmoError struct {
		ErrorDetail string `json:"error"`
		ErrorCode   int    `json:"error_code,string"`
	}
)
