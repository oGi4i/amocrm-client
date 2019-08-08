package amocrm

const (
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

	LeadCustomFieldInvalidId = "leads_custom_field_invalid_id"

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
		101: AccountNotFound,
		102: BodyMustBeJSON,
		103: InvalidRequestParameters,
		104: InvalidRequestMethod,

		110: InvalidCredentials,
		111: CaptchaInputRequired,
		112: DisabledAccount,
		113: UnauthorizedIpAddress,

		201: ContactAddEmptyArray,
		202: ContactAddInsufficientAccessRights,
		203: ContactAddCustomFieldsInternalError,
		204: ContactAddCustomFieldNotFound,
		205: ContactAddNotProcessed,

		206: ContactsEmptyRequest,
		207: ContactsInvalidRequestMethod,

		208: ContactUpdateEmptyArray,
		209: ContactUpdatedRequiredParametersMissing,
		210: ContactUpdateCustomFieldsInternalError,
		211: ContactUpdateCustomFieldNotFound,
		212: ContactUpdateNotProcessed,

		213: LeadAddEmptyArray,
		214: LeadEmptyRequest,
		215: LeadInvalidRequestMethod,
		216: LeadUpdateEmptyArray,
		217: LeadUpdateRequiredParametersMissing,

		219: ContactGetSearchError,

		240: LeadCustomFieldInvalidId,

		330: TooManyLinkedEntities,

		400: InvalidRequest,
		401: AccountNotFoundOnServer,
		402: SubscriptionExpire,
		403: AccountBlocked,

		429: RateLimitExceeded,

		2002: NoContent,
	}
)

type (
	AmoError struct {
		ErrorDetail string `json:"error"`
		ErrorCode   int    `json:"error_code,string"`
	}
)
