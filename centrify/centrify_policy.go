package centrify

import (
	"errors"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// Policy - Encapsulates a single policy
type Policy struct {
	vaultObject
	apiGetPolicies    string //`json:"-"`
	apiSetPolicyLinks string

	Plink    *PolicyLink `json:"Plink,omitempty" schema:"plink,omitempty"`
	Path     string      `json:"Path,omitempty" schema:"path,omitempty"`
	Position int         `json:"-,omitempty" schema:"position,omitempty"`
	Settings *settings   `json:"Settings,omitempty" schema:"settings,omitempty"`
}

type settings struct {
	CentrifyServices       *centrifyServices       `json:"CentrifyServices,omitempty" schema:"centrify_services,omitempty"`              // Authentication -> Centrify Services
	CentrifyClient         *centrifyClient         `json:"CentrifyClient,omitempty" schema:"centrify_client,omitempty"`                  // Authentication -> Centrify Clients -> Login
	CentrifyCSSServer      *centrifyCSSServer      `json:"CentrifyCSSServer,omitempty" schema:"centrify_css_server,omitempty"`           // Authentication -> Centrify Server Suite Agents -> Linux, UNIX and Windows Servers
	CentrifyCSSWorkstation *centrifyCSSWorkstation `json:"CentrifyCSSWorkstation,omitempty" schema:"centrify_css_workstation,omitempty"` // Authentication -> Centrify Server Suite Agents -> Windows Workstations
	CentrifyCSSElevation   *centrifyCSSElevation   `json:"CentrifyCSSElevation,omitempty" schema:"centrify_css_elevation,omitempty"`     // Authentication -> Centrify Server Suite Agents -> Privilege Elevation
	SelfService            *selfService            `json:"SelfService,omitempty" schema:"self_service,omitempty"`                        // User Security -> Self Service
	PasswordSettings       *passwordSettings       `json:"PasswordSettings,omitempty" schema:"password_settings,omitempty"`              // User Security -> Password Settings
	OATHOTP                *oathOTP                `json:"OATHOTP,omitempty" schema:"oath_otp,omitempty"`                                // User Security -> OATH OTP
	Radius                 *radius                 `json:"Radius,omitempty" schema:"radius,omitempty"`                                   // User Security -> RADIUS
	UserAccount            *userAccount            `json:"UserAccount,omitempty" schema:"user_account,omitempty"`                        // User Security -> User Account
	SystemSet              *systemSet              `json:"SystemSet,omitempty" schema:"system_set,omitempty"`                            // Resouces -> Systems
	DatabaseSet            *databaseSet            `json:"DatabaseSet,omitempty" schema:"database_set,omitempty"`                        // Resouces -> Databases
	DomainSet              *domainSet              `json:"DomainSet,omitempty" schema:"domain_set,omitempty"`                            // Resouces -> Domains
	AccountSet             *accountSet             `json:"AccountSet,omitempty" schema:"account_set,omitempty"`                          // Resouces -> Accounts
	SecretSet              *secretSet              `json:"SecretSet,omitempty" schema:"secret_set,omitempty"`                            // Resouces -> Secrets
	SSHKeySet              *sshKeySet              `json:"SSHKeySet,omitempty" schema:"sshkey_set,omitempty"`                            // Resouces -> SSH Keys
	MobileDevice           *mobileDevice           `json:"MobileDevice,omitempty" schema:"mobile_device,omitempty"`                      // Devices
}

// Authentication -> Centrify Services menu
// Authentication Policy for Centrify Services
type centrifyServices struct {
	// Session Parameters
	AuthenticationEnabled  bool            `json:"AuthenticationEnabled,omitempty" schema:"authentication_enabled,omitempty"`                                // Enable authentication policy controls
	DefaultProfileID       string          `json:"/Core/Authentication/AuthenticationRulesDefaultProfileId,omitempty" schema:"default_profile_id,omitempty"` // Default Profile (used if no conditions matched)
	ChallengeRules         *ChallengeRules `json:"/Core/Authentication/AuthenticationRules,omitempty" schema:"challenge_rule,omitempty"`
	SessionLifespan        int             `json:"/Core/Authentication/CookieSessionLifespanHours,omitempty" schema:"session_lifespan,omitempty"`         // Hours until session expires (default 12)
	AllowSessionPersist    bool            `json:"/Core/Authentication/CookieAllowPersist" schema:"allow_session_persist"`                                // Allow 'Keep me signed in' checkbox option at login (session spans browser sessions)
	DefaultSessionPersist  bool            `json:"/Core/Authentication/CookiePersistDefault,omitempty" schema:"default_session_persist,omitempty"`        // Default 'Keep me signed in' checkbox option to enabled
	PersistSessionLifespan int             `json:"/Core/Authentication/CookiePersistLifespanHours,omitempty" schema:"persist_session_lifespan,omitempty"` // Hours until session expires when 'Keep me signed in' option enabled (default 2 weeks)
	// Other Settings
	AllowIwa                   bool `json:"/Core/Authentication/AllowIwa" schema:"allow_iwa"`                                                             // Allow IWA connections (bypasses authentication rules and default profile)
	IwaSetKnownEndpoint        bool `json:"/Core/Authentication/IwaSetKnownEndpoint,omitempty" schema:"iwa_set_cookie,omitempty"`                         // Set identity cookie for IWA connections
	IwaSatisfiesAll            bool `json:"/Core/Authentication/IwaSatisfiesAllMechs,omitempty" schema:"iwa_satisfies_all,omitempty"`                     // IWA connections satisfy all MFA mechanisms
	UseCertAuth                bool `json:"/Core/Authentication/AllowZso" schema:"use_certauth"`                                                          // Use certificates for authentication
	CertAuthSkipChallenge      bool `json:"/Core/Authentication/ZsoSkipChallenge,omitempty" schema:"certauth_skip_challenge,omitempty"`                   // Certificate authentication bypasses authentication rules and default profile
	CertAuthSetKnownEndpoint   bool `json:"/Core/Authentication/ZsoSetKnownEndpoint,omitempty" schema:"certauth_set_cookie,omitempty"`                    // Set identity cookie for connections using certificate authentication
	CertAuthSatisfiesAll       bool `json:"/Core/Authentication/ZsoSatisfiesAllMechs,omitempty" schema:"certauth_satisfies_all,omitempty"`                // Connections using certificate authentication satisfy all MFA mechanisms
	NoMfaMechLogin             bool `json:"/Core/Authentication/NoMfaMechLogin" schema:"allow_no_mfa_mech"`                                               // Allow users without a valid authentication factor to log in
	FederatedLoginAllowsMfa    bool `json:"/Core/Authentication/FederatedLoginAllowsMfa" schema:"auth_rule_federated"`                                    // Apply additional authentication rules to federated users
	FederatedLoginSatisfiesAll bool `json:"/Core/Authentication/FederatedLoginSatisfiesAllMechs" schema:"federated_satisfies_all"`                        // Connections via Federation satisfy all MFA mechanisms
	BlockMechsOnMobileLogin    bool `json:"/Core/MfaRestrictions/BlockMobileMechsOnMobileLogin,omitempty" schema:"block_auth_from_same_device,omitempty"` // Allow additional authentication from same device
	ContinueFailedSessions     bool `json:"/Core/Authentication/ContinueFailedSessions" schema:"continue_failed_sessions"`                                // Continue with additional challenges after failed challenge
	SkipMechsInFalseAdvance    bool `json:"/Core/Authentication/SkipMechsInFalseAdvance,omitempty" schema:"stop_auth_on_prev_failed,omitempty"`           // Do not send challenge request when previous challenge response failed
	RememberLastAuthFactor     bool `json:"/Core/Authentication/AllowLoginMfaCache" schema:"remember_last_factor"`                                        // Remember and suggest last used authentication factor
}

type centrifyClient struct {
	AuthenticationEnabled bool            `json:"/Core/__centrify_cagent/AuthenticationEnabled,omitempty" schema:"authentication_enabled,omitempty"`                          // Enable authentication policy controls
	DefaultProfileID      string          `json:"/Core/__centrify_cagent/Authentication/AuthenticationRulesDefaultProfileId,omitempty" schema:"default_profile_id,omitempty"` // Default Profile (used if no conditions matched)
	ChallengeRules        *ChallengeRules `json:"/Core/__centrify_cagent/Authentication/AuthenticationRules,omitempty" schema:"challenge_rule,omitempty"`
	NoMfaMechLogin        bool            `json:"/Core/__centrify_cagent/Authentication/NoMfaMechLogin,omitempty" schema:"allow_no_mfa_mech,omitempty"` // Allow users without a valid authentication factor to log in
}

type centrifyCSSServer struct {
	AuthenticationEnabled bool            `json:"/Core/Css/AuthenticationEnabled,omitempty" schema:"authentication_enabled,omitempty"`                    // Enable authentication policy controls
	DefaultProfileID      string          `json:"/Core/Css/MfaLogin/AuthenticationRulesDefaultProfileId,omitempty" schema:"default_profile_id,omitempty"` // Default Profile (used if no conditions matched)
	ChallengeRules        *ChallengeRules `json:"/Core/Css/MfaLogin/AuthenticationRules,omitempty" schema:"challenge_rule,omitempty"`
	PassThroughMode       int             `json:"/Core/Css/MfaLogin/CssPinningMode,omitempty" schema:"pass_through_mode,omitempty"` // Apply pass-through duration
}

type centrifyCSSWorkstation struct {
	AuthenticationEnabled bool            `json:"/Core/Css/WindowsEndpointAuthenticationEnabled,omitempty" schema:"authentication_enabled,omitempty"`      // Enable authentication policy controls
	DefaultProfileID      string          `json:"/Core/Css/WinClient/AuthenticationRulesDefaultProfileId,omitempty" schema:"default_profile_id,omitempty"` // Default Profile (used if no conditions matched)
	ChallengeRules        *ChallengeRules `json:"/Core/Css/WinClient/AuthenticationRules,omitempty" schema:"challenge_rule,omitempty"`
}

type centrifyCSSElevation struct {
	AuthenticationEnabled bool            `json:"/Core/Css/PrivilegeElevationEnabled,omitempty" schema:"authentication_enabled,omitempty"`            // Enable authentication policy controls
	DefaultProfileID      string          `json:"/Core/Css/Dzdo/AuthenticationRulesDefaultProfileId,omitempty" schema:"default_profile_id,omitempty"` // Default Profile (used if no conditions matched)
	ChallengeRules        *ChallengeRules `json:"/Core/Css/Dzdo/AuthenticationRules,omitempty" schema:"challenge_rule,omitempty"`
}

type selfService struct {
	AccountSelfServiceEnabled bool `json:"PasswordResetEnabled,omitempty" schema:"account_selfservice_enabled,omitempty"` // Enable account self service controls
	// Password Reset
	PasswordResetEnabled         bool   `json:"/Core/PasswordReset/PasswordResetEnabled,omitempty" schema:"password_reset_enabled,omitempty"`
	PasswordResetADEnabled       bool   `json:"/Core/PasswordReset/PasswordResetADEnabled,omitempty" schema:"pwreset_allow_for_aduser,omitempty"`          // Allow for Active Directory users
	PasswordResetCookieOnly      bool   `json:"/Core/PasswordReset/PasswordResetIdentityCookieOnly,omitempty" schema:"pwreset_with_cookie_only,omitempty"` // Only allow from browsers with identity cookie
	PasswordResetRequiresRelogin bool   `json:"/Core/PasswordReset/PasswordResetRequiresMfaRestart,omitempty" schema:"login_after_reset,omitempty"`        // User must log in after successful password reset
	PasswordResetAuthProfile     string `json:"/Core/PasswordReset/PasswordResetAuthProfile,omitempty" schema:"pwreset_auth_profile_id,omitempty"`         // Password reset authentication profile
	PasswordResetMaxAttempts     int    `json:"/Core/PasswordReset/PasswordResetMaxAttemptsPerSession,omitempty" schema:"max_reset_attempts,omitempty"`    // Maximum consecutive password reset attempts per session
	// Account Unlock
	AccountUnlockEnabled     bool   `json:"/Core/PasswordReset/AccountUnlockEnabled,omitempty" schema:"account_unlock_enabled,omitempty"`
	AccountUnlockADEnabled   bool   `json:"/Core/PasswordReset/AccountUnlockADEnabled,omitempty" schema:"unlock_allow_for_aduser,omitempty"`          // Allow for Active Directory users
	AccountUnlockCookieOnly  bool   `json:"/Core/PasswordReset/AccountUnlockIdentityCookieOnly,omitempty" schema:"unlock_with_cookie_only,omitempty"` // Only allow from browsers with identity cookie
	ShowAccountLocked        bool   `json:"/Mobile/EndpointAgent/showAccountLocked" schema:"show_locked_message,omitempty"`                           // Show a message to end users in desktop login that account is locked (default no)
	AccountUnlockAuthProfile string `json:"/Core/PasswordReset/AccountUnlockAuthProfile,omitempty" schema:"unlock_auth_profile_id,omitempty"`         // Account unlock authentication profile
	// Active Directory Self Service Settings
	UseADAdmin  bool         `json:"/Core/PasswordReset/UseADAdmin,omitempty" schema:"use_ad_admin,omitempty"` // Use AD admin for AD self-service
	ADAdminUser string       `json:"/Core/PasswordReset/ADAdminUser,omitempty" schema:"ad_admin_user,omitempty"`
	ADAdminPass *adAdminPass `json:"/Core/PasswordReset/ADAdminPass,omitempty" schema:"admin_user_password,omitempty"`
	// Additional Policy Parameters
	MaxResetAllowed int `json:"/Core/PasswordReset/Max,omitempty" schema:"max_reset_allowed,omitempty"`    // Maximum forgotten password resets allowed within window (default 10)
	MaxTimeAllowed  int `json:"/Core/PasswordReset/MaxTime,omitempty" schema:"max_time_allowed,omitempty"` // Capture window for forgotten password resets (default 60 minutes)
}

type adAdminPass struct {
	Type  string `json:"_Type,omitempty" schema:"type,omitempty"`
	Value string `json:"_Value,omitempty" schema:"value,omitempty"`
}

type passwordSettings struct {
	// Password Requirements
	MinLength      int  `json:"/Core/Security/CDS/PasswordPolicy/MinLength,omitempty" schema:"min_length,omitempty"`            // Minimum password length (default 8)
	MaxLength      int  `json:"/Core/Security/CDS/PasswordPolicy/MaxLength,omitempty" schema:"max_length,omitempty"`            // Maximum password length (default 64)
	RequireDigit   bool `json:"/Core/Security/CDS/PasswordPolicy/RequireDigit,omitempty" schema:"require_digit,omitempty"`      // Require at least one digit (default yes)
	RequireMixCase bool `json:"/Core/Security/CDS/PasswordPolicy/RequireMixCase,omitempty" schema:"require_mix_case,omitempty"` // Require at least one upper case and one lower case letter (default yes)
	RequireSymbol  bool `json:"/Core/Security/CDS/PasswordPolicy/RequireSymbol,omitempty" schema:"require_symbol,omitempty"`    // Require at least one symbol (default no)
	// Display Requirements
	ShowPasswordComplexity bool   `json:"/Core/Security/CDS/PasswordPolicy/ShowPasswordComplexity,omitempty" schema:"show_password_complexity,omitempty"` // Show password complexity requirements when entering a new password (default no)
	NonCdsComplexityHint   string `json:"/Core/Security/CDS/PasswordPolicy/NonCdsComplexityHint,omitempty" schema:"complexity_hint,omitempty"`            // Password complexity requirements for directory services other than Centrify Directory
	// Additional Requirements
	AllowRepeatedChar       int  `json:"/Core/Security/CDS/PasswordPolicy/AllowRepeatedChar,omitempty" schema:"no_of_repeated_char_allowed,omitempty"`     // Limit the number of consecutive repeated characters
	CheckWeakPassword       bool `json:"/Core/Security/CDS/PasswordPolicy/CheckWeakPassword,omitempty" schema:"check_weak_password,omitempty"`             // Check against weak password
	AllowIncludeUsername    bool `json:"/Core/Security/CDS/PasswordPolicy/AllowIncludeUsername,omitempty" schema:"allow_include_username,omitempty"`       // Allow username as part of password
	AllowIncludeDisplayname bool `json:"/Core/Security/CDS/PasswordPolicy/AllowIncludeDisplayname,omitempty" schema:"allow_include_displayname,omitempty"` // Allow display name as part of password
	RequireUnicode          bool `json:"/Core/Security/CDS/PasswordPolicy/RequireUnicode,omitempty" schema:"require_unicode,omitempty"`                    // Require at least one Unicode characters
	// Password Age
	MinAgeInDays   int  `json:"/Core/Security/CDS/PasswordPolicy/MinAgeInDays,omitempty" schema:"min_age_in_days,omitempty"` // Minimum password age before change is allowed (default 0 days)
	MaxAgeInDays   int  `json:"/Core/Security/CDS/PasswordPolicy/AgeInDays,omitempty" schema:"max_age_in_days,omitempty"`    // Maximum password age (default 365 days)
	History        int  `json:"/Core/Security/CDS/PasswordPolicy/History,omitempty" schema:"password_history,omitempty"`     // Password history (default 3)
	NotifySoft     int  `json:"/Core/PasswordReset/NotifySoft,omitempty" schema:"expire_soft_notification,omitempty"`        // Password Expiration Notification (default 14 days)
	NotifyHard     int  `json:"/Core/PasswordReset/NotifyHard,omitempty" schema:"expire_hard_notification,omitempty"`        // Escalated Password Expiration Notification (default 48 hours)
	NotifyOnMobile bool `json:"/Core/PasswordChange/NotifyOnMobile,omitempty" schema:"expire_notification_mobile,omitempty"` // Enable password expiration notifications on enrolled mobile devices
	// Capture Settings
	BadAttemptThreshold int `json:"/Core/Security/CDS/LockoutPolicy/Threshold,omitempty" schema:"bad_attempt_threshold,omitempty"` // Maximum consecutive bad password attempts allowed within window (default Off)
	CaptureWindow       int `json:"/Core/Security/CDS/LockoutPolicy/Window,omitempty" schema:"capture_window,omitempty"`           // Capture window for consecutive bad password attempts (default 30 minutes)
	LockoutDuration     int `json:"/Core/Security/CDS/LockoutPolicy/Duration,omitempty" schema:"lockout_duration,omitempty"`       // Lockout duration before password re-attempt allowed (default 30 minutes)
}

type oathOTP struct {
	AllowOTP bool `json:"/Core/Security/CDS/ExternalMFA/ShowQRCode,omitempty" schema:"allow_otp,omitempty"` // Allow OATH OTP integration
}

type radius struct {
	AllowRadius          bool   `json:"/Core/Authentication/AllowRadius,omitempty" schema:"allow_radius,omitempty"`                                  // Allow RADIUS client connections
	RadiusUseChallenges  bool   `json:"/Core/Authentication/RadiusUseChallenges,omitempty" schema:"require_challenges,omitempty"`                    // Require authentication challenge
	DefaultProfileID     string `json:"/Core/Authentication/RadiusChallengeProfile,omitempty" schema:"default_profile_id,omitempty"`                 // Default authentication profile
	SendVendorAttributes bool   `json:"/Core/Authentication/SendRadiusVendorSpecificAttributes,omitempty" schema:"send_vendor_attributes,omitempty"` // Send vendor specific attributes
	AllowExternalRadius  bool   `json:"/Core/Authentication/AllowExternalRadius,omitempty" schema:"allow_external_radius,omitempty"`                 // Allow 3rd Party RADIUS Authentication
}

type userAccount struct {
	UserChangePasswordAllow     bool   `json:"/Core/PasswordChange/UserChangeAllow,omitempty" schema:"allow_user_change_password,omitempty"`                  // Enable users to change their passwords
	PasswordChangeAuthProfileID string `json:"/Core/Authentication/UserUpdateProfile/Password,omitempty" schema:"password_change_auth_profile_id,omitempty"`  // Authentication Profile required to change password
	ShowU2f                     bool   `json:"/Core/Security/CDS/ExternalMFA/ShowU2f,omitempty" schema:"show_fido2,omitempty"`                                // Enable users to enroll FIDO2 Authenticators
	U2fPrompt                   string `json:"/Core/Security/CDS/ExternalMFA/U2fUiPrompt,omitempty" schema:"fido2_prompt,omitempty"`                          // FIDO2 Security Key Display Name
	U2fAuthProfileID            string `json:"/Core/Authentication/UserUpdateProfile/U2F,omitempty" schema:"fido2_auth_profile_id,omitempty"`                 // Authentication Profile required to configure FIDO2 Authenticators
	ShowQRCode                  bool   `json:"/Core/Security/CDS/ExternalMFA/ShowQRCodeForSelfService,omitempty" schema:"show_otp,omitempty"`                 // Enable users to configure an OATH OTP client (requires enabling OATH OTP policy)
	OTPPrompt                   string `json:"/Core/Security/CDS/ExternalMFA/UiPrompt,omitempty" schema:"otp_prompt,omitempty"`                               // OATH OTP Display Name
	OTPAuthProfileID            string `json:"/Core/Authentication/UserUpdateProfile/OathProfile,omitempty" schema:"otp_auth_profile_id,omitempty"`           // Authentication Profile required to configure OATH OTP client
	ConfigureSecurityQuestions  bool   `json:"/Core/Authentication/ConfigureSecurityQuestions,omitempty" schema:"configure_security_questions,omitempty"`     // Enable users to configure Security Questions
	AllowDupAnswers             bool   `json:"/Core/Authentication/SecurityQuestionPreventDupAnswers,omitempty" schema:"prevent_dup_answers,omitempty"`       // Allow duplicate security question answers
	UserDefinedQuestions        int    `json:"/Core/Authentication/UserSecurityQuestionsPerUser,omitempty" schema:"user_defined_questions,omitempty"`         // Required number of user-defined questions
	AdminDefinedQuestions       int    `json:"/Core/Authentication/AdminSecurityQuestionsPerUser,omitempty" schema:"admin_defined_questions,omitempty"`       // Required number of admin-defined questions
	MinCharInAnswer             int    `json:"/Core/Authentication/SecurityQuestionAnswerMinLength,omitempty" schema:"min_char_in_answer,omitempty"`          // Minimum number of characters required in answers
	QuestionAuthProfileID       string `json:"/Core/Authentication/UserUpdateProfile/SecurityQuestion,omitempty" schema:"question_auth_profile_id,omitempty"` // Authentication Profile required to set Security Questions
	PhonePinChangeAllow         bool   `json:"/Core/PhoneAuth/UserChangeAllow,omitempty" schema:"allow_phone_pin_change,omitempty"`                           // Enable users to configure a Phone PIN for MFA
	MinPhonePinLength           int    `json:"/Core/Authentication/MinPhonePinLength,omitempty" schema:"min_phone_pin_length,omitempty"`                      // Minimum Phone PIN length
	PhonePinAuthProfileID       string `json:"/Core/Authentication/UserUpdateProfile/PhonePin,omitempty" schema:"phone_pin_auth_profile_id,omitempty"`        // Authentication Profile required to configure a Phone PIN
	AllowUserChangeMFARedirect  bool   `json:"/Core/Security/CDS/AllowUserChangeMFARedirect,omitempty" schema:"allow_mfa_redirect_change,omitempty"`          // Enable users to redirect multi factor authentication to a different user account
	UserProfileAuthProfileID    string `json:"/Core/Authentication/UserUpdateProfile/Profile,omitempty" schema:"user_profile_auth_profile_id,omitempty"`      // Authentication Profile required to modify Personal Profile
	DefaultLanguage             string `json:"/Core/Policy/Culture,omitempty" schema:"default_language,omitempty"`                                            // Default Language
}

type systemSet struct {
	// Account Policy
	DefaultCheckoutTime int `json:"/PAS/Server/DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	// System Policy
	AllowRemote         bool            `json:"/PAS/Server/AllowRemote,omitempty" schema:"allow_remote_access,omitempty"`        // Allow access from a public network (web client only)
	AllowRdpClipboard   bool            `json:"/PAS/Server/AllowRdpClipboard,omitempty" schema:"allow_rdp_clipboard,omitempty"`  // Allow RDP client to sync local clipboard with remote session
	LoginDefaultProfile string          `json:"/PAS/Server/LoginDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default System Login Profile (used if no conditions matched)
	ChallengeRules      *ChallengeRules `json:"/PAS/Server/LoginRules,omitempty" schema:"challenge_rule,omitempty"`
	// Security Settings
	AllowMultipleCheckouts            bool   `json:"/PAS/ConfigurationSetting/Server/AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts for this system
	AllowPasswordRotation             bool   `json:"/PAS/ConfigurationSetting/Server/AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int    `json:"/PAS/ConfigurationSetting/Server/PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool   `json:"/PAS/ConfigurationSetting/Server/AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int    `json:"/PAS/ConfigurationSetting/Server/MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	MinimumSSHKeysAge                 int    `json:"/PAS/ConfigurationSetting/Server/MinimumSshKeysAge,omitempty" schema:"minimum_sshkey_age,omitempty"`                                     // Minimum SSH Key Age (days)
	AllowSSHKeysRotation              bool   `json:"/PAS/ConfigurationSetting/Server/AllowSshKeysRotation,omitempty" schema:"enable_sshkey_rotation,omitempty"`                              // Enable periodic SSH key rotation
	SSHKeysRotateDuration             int    `json:"/PAS/ConfigurationSetting/Server/SshKeysRotateDuration,omitempty" schema:"sshkey_rotate_interval,omitempty"`                             // SSH key rotation interval (days)
	SSHKeysGenerationAlgorithm        string `json:"/PAS/ConfigurationSetting/Server/SshKeysGenerationAlgorithm,omitempty" schema:"sshkey_algorithm,omitempty"`                              // SSH Key Generation Algorithm
	// Maintenance Settings
	AllowPasswordHistoryCleanUp    bool `json:"/PAS/ConfigurationSetting/Server/AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`     // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration int  `json:"/PAS/ConfigurationSetting/Server/PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"` // Password history cleanup (days)
	AllowSSHKeysCleanUp            bool `json:"/PAS/ConfigurationSetting/Server/AllowSshKeysCleanUp,omitempty" schema:"enable_sshkey_history_cleanup,omitempty"`               // Enable periodic SSH key cleanup
	SSHKeysCleanUpDuration         int  `json:"/PAS/ConfigurationSetting/Server/SshKeysCleanUpDuration,omitempty" schema:"sshkey_historycleanup_duration,omitempty"`           // SSH key cleanup (days)
}

type databaseSet struct {
	// Account Policy
	DefaultCheckoutTime int `json:"/PAS/VaultDatabase/DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	// Security Settings
	AllowMultipleCheckouts            bool `json:"/PAS/ConfigurationSetting/VaultDatabase/AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts for related accounts
	AllowPasswordRotation             bool `json:"/PAS/ConfigurationSetting/VaultDatabase/AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int  `json:"/PAS/ConfigurationSetting/VaultDatabase/PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool `json:"/PAS/ConfigurationSetting/VaultDatabase/AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int  `json:"/PAS/ConfigurationSetting/VaultDatabase/MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	// Maintenance Settings
	AllowPasswordHistoryCleanUp    bool `json:"/PAS/ConfigurationSetting/VaultDatabase/AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`     // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration int  `json:"/PAS/ConfigurationSetting/VaultDatabase/PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"` // Password history cleanup (days)
}

type domainSet struct {
	// Account Policy
	DefaultCheckoutTime int `json:"/PAS/VaultDomain/DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	// Security Settings
	AllowMultipleCheckouts            bool `json:"/PAS/ConfigurationSetting/VaultDomain/AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts per AD account added for this domain
	AllowPasswordRotation             bool `json:"/PAS/ConfigurationSetting/VaultDomain/AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int  `json:"/PAS/ConfigurationSetting/VaultDomain/PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool `json:"/PAS/ConfigurationSetting/VaultDomain/AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int  `json:"/PAS/ConfigurationSetting/VaultDomain/MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	// Maintenance Settings
	AllowPasswordHistoryCleanUp    bool `json:"/PAS/ConfigurationSetting/VaultDomain/AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`     // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration int  `json:"/PAS/ConfigurationSetting/VaultDomain/PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"` // Password history cleanup (days)
}

type accountSet struct {
	// Account Security
	DefaultCheckoutTime            int             `json:"/PAS/VaultAccount/DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"`   // Checkout lifetime (minutes)
	PasswordCheckoutDefaultProfile string          `json:"/PAS/VaultAccount/PasswordCheckoutDefaultProfile" schema:"default_profile_id,omitempty"` // Default Password Checkout Profile (used if no conditions matched)
	ChallengeRules                 *ChallengeRules `json:"/PAS/VaultAccount/PasswordCheckoutRules,omitempty" schema:"challenge_rule,omitempty"`
}

type secretSet struct {
	DataVaultDefaultProfile string          `json:"/PAS/DataVault/DataVaultDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default Secret Challenge Profile (used if no conditions matched)
	ChallengeRules          *ChallengeRules `json:"/PAS/DataVault/DataVaultRules,omitempty" schema:"challenge_rule,omitempty"`
}

type sshKeySet struct {
	SSHKeysDefaultProfile string          `json:"/PAS/SshKeys/SshKeysDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default SSH Key Challenge Profile
	ChallengeRules        *ChallengeRules `json:"/PAS/SshKeys/SshKeysRules,omitempty" schema:"challenge_rule,omitempty"`
}

type mobileDevice struct {
	AllowEnrollment           bool `json:"/Mobile/EnrollRules/Common/AllowEnrollment,omitempty" schema:"allow_enrollment,omitempty"`                                                                                              // Permit device registration
	AllowJailBrokenDevices    bool `json:"/Mobile/EnrollRules/Common/AllowJailBrokenDevices,omitempty" schema:"permit_non_compliant_device,omitempty"`                                                                            // Permit non-compliant devices to register
	EnableInviteEnrollment    bool `json:"/Mobile/DeviceManagement/EnableInviteBasedEnrollment,omitempty" schema:"enable_invite_enrollment,omitempty"`                                                                            // Enable invite based registration
	AllowNotifnMutipleDevices bool `json:"/Mobile/Software/Policies/Centrify/Common/AllowNotificationOnMutipleDevices,omitempty" schema:"allow_notify_multi_devices,omitempty"`                                                   // Allow user notifications on multiple devices
	AllowDebugLogging         bool `json:"/Mobile/Software/Policies/Centrify/iOSSettings/AllowDebugLogging,omitempty" schema:"enable_debug,omitempty"`                                                                            // Enable debug logging
	LocationTracking          bool `json:"/Mobile/Software/Policies/Centrify/iOSSettings/Restrictions/LocationTracking,omitempty" schema:"location_tracking,omitempty"`                                                           // Report mobile device location
	ForceFingerprint          bool `json:"/Mobile/Software/Policies/Centrify/Application/Security/MobileAuthenticator/ForceFingerprintForMobileAuthenticator,omitempty" schema:"force_fingerprint,omitempty"`                     // Enforce fingerprint scan for Mobile Authenticator
	AllowFallbackAppPin       bool `json:"/Mobile/Software/Policies/Centrify/Application/Security/MobileAuthenticator/ForceFingerprintForMobileAuthenticatorAllowFallbackAppPin,omitempty" schema:"allow_fallback_pin,omitempty"` // Allow App PIN
	RequestPasscode           bool `json:"/Mobile/Software/Policies/Centrify/Application/Passcode/ForceAppPin,omitempty" schema:"require_passcode,omitempty"`                                                                     // Require client application passcode on device
	AutoLockTimeout           int  `json:"/Mobile/Software/Policies/Centrify/Application/Passcode/AppInactivityTimeout,omitempty" schema:"auto_lock_timeout,omitempty"`                                                           // Auto-Lock (minutes)
	AppLockOnExit             bool `json:"/Mobile/Software/Policies/Centrify/Application/Passcode/AppLockOnExit,omitempty" schema:"lock_app_on_exit,omitempty"`                                                                   // Lock on exit
}

// NewPolicy is a policy constructor
func NewPolicy(c *restapi.RestClient) *Policy {
	s := Policy{}
	s.Plink = &PolicyLink{}
	s.client = c
	s.apiRead = "/Policy/GetPolicyBlock"
	s.apiCreate = "/Policy/SavePolicyBlock3"
	s.apiDelete = "/Policy/DeletePolicyBlock"
	s.apiUpdate = "/Policy/SavePolicyBlock3"
	s.apiGetPolicies = "/Policy/GetNicePlinks"
	s.apiSetPolicyLinks = "/Policy/setPlinksv2"

	return &s
}

// Read function fetches a Policy from source, including attribute values. Returns error if any
func (o *Policy) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}

	var queryArg = make(map[string]interface{})
	// We assume name is always the same as ID. In actual fact, this isn't the case. But we won't support policy rename so this is ok
	queryArg["name"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	//LogD.Printf("Response for Policy from tenant: %+v", resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}

	// Fill root level attributes: Path, Description
	fillWithMap(o, resp.Result)

	// Fill root level attributes: Params, LinkType, PolicySet, Description
	resp2, err2 := o.Query("")
	LogD.Printf("Response for Policy query: %+v", resp2)
	if err2 != nil {
		return err2
	}

	plink := &PolicyLink{}
	fillWithMap(plink, resp2)
	o.Plink = plink
	// Fill settings
	var settings = resp.Result["Settings"].(map[string]interface{})

	CentrifyServices := &centrifyServices{}
	fillWithMap(CentrifyServices, settings)
	o.Settings.CentrifyServices = CentrifyServices

	CentrifyClient := &centrifyClient{}
	fillWithMap(CentrifyClient, settings)
	o.Settings.CentrifyClient = CentrifyClient

	CentrifyCSSServer := &centrifyCSSServer{}
	fillWithMap(CentrifyCSSServer, settings)
	o.Settings.CentrifyCSSServer = CentrifyCSSServer

	CentrifyCSSWorkstation := &centrifyCSSWorkstation{}
	fillWithMap(CentrifyCSSWorkstation, settings)
	o.Settings.CentrifyCSSWorkstation = CentrifyCSSWorkstation

	CentrifyCSSElevation := &centrifyCSSElevation{}
	fillWithMap(CentrifyCSSElevation, settings)
	o.Settings.CentrifyCSSElevation = CentrifyCSSElevation

	SelfService := &selfService{}
	fillWithMap(SelfService, settings)
	o.Settings.SelfService = SelfService

	PasswordSettings := &passwordSettings{}
	fillWithMap(PasswordSettings, settings)
	o.Settings.PasswordSettings = PasswordSettings

	OATHOTP := &oathOTP{}
	fillWithMap(OATHOTP, settings)
	o.Settings.OATHOTP = OATHOTP

	Radius := &radius{}
	fillWithMap(Radius, settings)
	o.Settings.Radius = Radius

	UserAccount := &userAccount{}
	fillWithMap(UserAccount, settings)
	o.Settings.UserAccount = UserAccount

	SystemSet := &systemSet{}
	fillWithMap(SystemSet, settings)
	o.Settings.SystemSet = SystemSet

	DatabaseSet := &databaseSet{}
	fillWithMap(DatabaseSet, settings)
	o.Settings.DatabaseSet = DatabaseSet

	DomainSet := &domainSet{}
	fillWithMap(DomainSet, settings)
	o.Settings.DomainSet = DomainSet

	AccountSet := &accountSet{}
	fillWithMap(AccountSet, settings)
	o.Settings.AccountSet = AccountSet

	SecretSet := &secretSet{}
	fillWithMap(SecretSet, settings)
	o.Settings.SecretSet = SecretSet

	SSHKeySet := &sshKeySet{}
	fillWithMap(SSHKeySet, settings)
	o.Settings.SSHKeySet = SSHKeySet

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Delete function deletes a Policy and returns a map that contains deletion result
func (o *Policy) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("path")
}

// Create function creates a Policy and returns a map that contains update result
func (o *Policy) Create() (*restapi.GenericMapResponse, error) {
	var queryArg = make(map[string]interface{})

	// Handle plinks
	plinks, _, err := o.getPlinks()
	if err != nil {
		return nil, err
	}
	var plink = make(map[string]interface{})
	plink["Name"] = o.Name
	plink["Description"] = o.Description
	plink["LinkType"] = o.Plink.LinkType
	plink["PolicySet"] = "/Policy/" + o.Name
	plink["Params"] = o.Plink.Params
	plink["Priority"] = 1

	plinks = insert(plinks, plink, 0)
	queryArg["plinks"] = plinks

	// Convert to nested map
	nestedmap, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	//LogD.Printf("Nested map for Create(): %+v", nestedmap)

	// Flatten Settings
	var settings = make(map[string]interface{})
	//flattenNestedMap(settings, nestedmap["Settings"])
	flattenSettings(settings, nestedmap["Settings"])
	// Remove Settings key which is nested map and replace it with flattened map
	delete(nestedmap, "Settings")
	policy := nestedmap
	policy["Path"] = "/Policy/" + o.Name
	policy["Settings"] = settings
	policy["Newpolicy"] = true
	queryArg["policy"] = policy

	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	reply, err := o.client.CallGenericMapAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Update function updates an existing Policy and returns a map that contains update result
func (o *Policy) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})

	// Convert to nested map
	nestedmap, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

	// Handle policy link
	plinks, _ := o.constructPlinks()
	queryArg["plinks"] = plinks
	delete(nestedmap, "Plink")

	// Flatten Settings
	var settings = make(map[string]interface{})
	//flattenNestedMap(settings, nestedmap["Settings"])
	flattenSettings(settings, nestedmap["Settings"])
	// Remove Settings key which is nested map and replace it with flattened map
	delete(nestedmap, "Settings")

	policy := nestedmap
	policy["Path"] = "/Policy/" + o.Name
	policy["Settings"] = settings
	policy["Newpolicy"] = false

	// Read policy again to fetch its latest RevStamp
	var queryArg2 = make(map[string]interface{})
	queryArg2["name"] = o.ID
	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg2)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	policy["RevStamp"] = resp.Result["RevStamp"]

	queryArg["policy"] = policy

	LogD.Printf("Generated Map for Update(): %+v", queryArg)

	reply, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single Policy object in map format
func (o *Policy) Query(key string) (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	resp, err := o.client.CallGenericMapAPI(o.apiGetPolicies, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	// Loop through respond results
	var results = resp.Result["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		LogD.Printf("Query row: %+v", row)
		if strings.EqualFold(key, "name") {
			if row["PolicySet"] == "/Policy/"+o.Name {
				return row, nil
			}
		} else {
			if row["ID"] == o.ID {
				return row, nil
			}
		}
	}

	return nil, errors.New("Found 0 matched policy")
}

func (o *Policy) getPlinks() ([]map[string]interface{}, string, error) {
	var plinks []map[string]interface{}

	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	resp, err := o.client.CallGenericMapAPI(o.apiGetPolicies, queryArg)
	if err != nil {
		return nil, "", err
	}
	if !resp.Success {
		return nil, "", errors.New(resp.Message)
	}

	var rev = resp.Result["RevStamp"].(string)
	var results = resp.Result["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		plinks = append(plinks, row)
	}

	return plinks, rev, nil
}

// constructPlinks updates the attributes in plinks section of update request
func (o *Policy) constructPlinks() ([]map[string]interface{}, string) {
	var plinks []map[string]interface{}
	oldplinks, rev, err := o.getPlinks()
	if err != nil {
		return nil, ""
	}
	// Loop through exiting plinks, find matching plink and remove it
	for _, v := range oldplinks {
		if o.ID != v["ID"] {
			plinks = append(plinks, v)
		} else {
			var plink = make(map[string]interface{})
			plink["Name"] = o.Name
			plink["ID"] = o.ID
			plink["Description"] = o.Description
			plink["LinkType"] = o.Plink.LinkType
			plink["PolicySet"] = "/Policy/" + o.Name
			plink["Params"] = o.Plink.Params
			plinks = append(plinks, plink)
		}
	}

	return plinks, rev
}

func (o *Policy) validateSettings() error {
	if o.Settings != nil {
		if o.Settings.CentrifyServices != nil {
			if data := o.Settings.CentrifyServices; data != nil {
				if data.AuthenticationEnabled && data.DefaultProfileID == "" {
					return errors.New("In CentrifyServices: AuthenticationEnabled is true but DefaultProfileID is empty")
				}
				if data.FederatedLoginAllowsMfa && data.FederatedLoginSatisfiesAll {
					return errors.New("In CentrifyServices: FederatedLoginAllowsMfa & FederatedLoginSatisfiesAll only one should be enabled")
				}
			}
		}

		if o.Settings.CentrifyClient != nil {
			if data := o.Settings.CentrifyClient; data != nil {
				if data.AuthenticationEnabled && data.DefaultProfileID == "" {
					return errors.New("In CentrifyClient: AuthenticationEnabled is true but DefaultProfileID is empty")
				}
			}
		}

		if o.Settings.CentrifyCSSServer != nil {
			if data := o.Settings.CentrifyCSSServer; data != nil {
				if data.AuthenticationEnabled && data.DefaultProfileID == "" {
					return errors.New("In CentrifyCSSServer: AuthenticationEnabled is true but DefaultProfileID is empty")
				}
			}
		}

		if o.Settings.CentrifyCSSWorkstation != nil {
			if data := o.Settings.CentrifyCSSWorkstation; data != nil {
				if data.AuthenticationEnabled && data.DefaultProfileID == "" {
					return errors.New("In CentrifyCSSWorkstation: AuthenticationEnabled is true but DefaultProfileID is empty")
				}
			}
		}

		if o.Settings.CentrifyCSSElevation != nil {
			if data := o.Settings.CentrifyCSSElevation; data != nil {
				if data.AuthenticationEnabled && data.DefaultProfileID == "" {
					return errors.New("In CentrifyCSSElevation: AuthenticationEnabled is true but DefaultProfileID is empty")
				}
			}
		}

		if o.Settings.SelfService != nil {
			if data := o.Settings.SelfService; data != nil {
				if data.PasswordResetEnabled && data.PasswordResetAuthProfile == "" {
					return errors.New("In SelfService: PasswordResetEnabled is true but PasswordResetAuthProfile is empty")
				}
				if data.AccountUnlockEnabled && data.AccountUnlockAuthProfile == "" {
					return errors.New("In SelfService: AccountUnlockEnabled is true but AccountUnlockAuthProfile is empty")
				}
			}
		}

		if o.Settings.PasswordSettings != nil {
			data := o.Settings.PasswordSettings
			if data.MinLength > data.MaxLength {
				return errors.New("In PasswordSettings: MinLength must be smaller than MaxLength")
			}
			if data.NotifySoft > data.NotifyHard {
				return errors.New("In PasswordSettings: NotifySoft must be smaller than NotifyHard")
			}
		}

		if o.Settings.Radius != nil {
			data := o.Settings.Radius
			if !data.AllowRadius && (data.RadiusUseChallenges || data.SendVendorAttributes) {
				return errors.New("In Radius: AllowRadius must be enabled before RadiusUseChallenges or SendVendorAttributes can be enabled")
			}
			if data.RadiusUseChallenges && data.DefaultProfileID == "" {
				return errors.New("In Radius: RadiusUseChallenges is true but DefaultProfileID is empty")
			}
		}

		if o.Settings.UserAccount != nil {
			data := o.Settings.UserAccount
			if data.ShowU2f && data.U2fPrompt == "" {
				return errors.New("In SelfService: ShowU2f is true but U2fPrompt is empty")
			}
			if data.ShowQRCode && data.OTPPrompt == "" {
				return errors.New("In SelfService: ShowQRCode is true but OTPPrompt is empty")
			}
			if data.ConfigureSecurityQuestions && (data.UserDefinedQuestions == 0 || data.AdminDefinedQuestions == 0 || data.MinCharInAnswer == 0) {
				return errors.New("In SelfService: ConfigureSecurityQuestions is true but UserDefinedQuestions or AdminDefinedQuestions or MinCharInAnswer is 0")
			}
		}

		if o.Settings.SystemSet != nil {
			data := o.Settings.SystemSet
			if data.AllowPasswordRotation && data.PasswordRotateDuration == 0 {
				return errors.New("In SystemSet: AllowPasswordRotation is true but PasswordRotateDuration is empty")
			}
			if data.AllowSSHKeysRotation && data.SSHKeysRotateDuration == 0 {
				return errors.New("In SystemSet: AllowSSHKeysRotation is true but SSHKeysRotateDuration is empty")
			}
			if data.AllowPasswordHistoryCleanUp && data.PasswordHistoryCleanUpDuration == 0 {
				return errors.New("In SystemSet: AllowPasswordHistoryCleanUp is true but PasswordHistoryCleanUpDuration is empty")
			}
			if data.AllowSSHKeysCleanUp && data.SSHKeysCleanUpDuration == 0 {
				return errors.New("In SystemSet: AllowSSHKeysCleanUp is true but SSHKeysCleanUpDuration is empty")
			}
		}

		if o.Settings.DatabaseSet != nil {
			data := o.Settings.DatabaseSet
			if data.AllowPasswordRotation && data.PasswordRotateDuration == 0 {
				return errors.New("In DatabaseSet: AllowPasswordRotation is true but PasswordRotateDuration is empty")
			}
			if data.AllowPasswordHistoryCleanUp && data.PasswordHistoryCleanUpDuration == 0 {
				return errors.New("In DatabaseSet: AllowPasswordHistoryCleanUp is true but PasswordHistoryCleanUpDuration is empty")
			}
		}

		if o.Settings.DomainSet != nil {
			data := o.Settings.DomainSet
			if data.AllowPasswordRotation && data.PasswordRotateDuration == 0 {
				return errors.New("In DomainSet: AllowPasswordRotation is true but PasswordRotateDuration is empty")
			}
			if data.AllowPasswordHistoryCleanUp && data.PasswordHistoryCleanUpDuration == 0 {
				return errors.New("In DomainSet: AllowPasswordHistoryCleanUp is true but PasswordHistoryCleanUpDuration is empty")
			}
		}
	}
	return nil
}

/*
	API to manage policy


	Fetch policy
	https://developer.centrify.com/reference#post_policy-getpolicyblock

		Request body format
		{
			"name": "/Policy/LAB User Login Policy"
		}

		Respond result
        {
            "success": true,
            "Result": {
                "Version": 1,
                "Settings": {
                    "/Core/Authentication/AllowLoginMfaCache": false,
                    "/Core/PasswordChange/UserChangeAllow": true,
                    "/Core/Authentication/ContinueFailedSessions": true,
                    "/Core/Authentication/AuthenticationRules": {
                        "_UniqueKey": "Condition",
                        "_Value": [],
                        "Enabled": true,
                        "_Type": "RowSet"
                    },

                    "/Core/Authentication/FederatedLoginSatisfiesAllMechs": false,
                    "/Core/Security/CDS/ExternalMFA/UiPrompt": "OATH OTP Client",
                    "/Core/Authentication/AuthenticationRulesHighAuthRequestedProfileId": "AlwaysAllowed"
                },
                "RevStamp": "637303716720000000",
                "RadiusClientList": [],
                "AuthProfiles": [
                    {
                        "Uuid": "cf484a69-4f1b-42e0-9286-2e0ce972b362",
                        "Name": "Default New Device Login Profile",
                        "DurationInMinutes": 720,
                        "Challenges": [
                            "UP",
                            "OTP,SMS,EMAIL,OATH,SQ"
                        ],
                        "AdditionalData": {
                            "NumberOfQuestions": 1
                        }
                    },
                    {
                        "Uuid": "08ccf738-9277-479e-83f4-374609dd9f0e",
                        "Name": "Default Other Login Profile",
                        "DurationInMinutes": 720,
                        "Challenges": [
                            "UP"
                        ],
                        "AdditionalData": {}
                    },
                    {
                        "Uuid": "624a5e94-da4c-4ba3-a596-f56b60637976",
                        "Name": "Default Password Reset Profile",
                        "DurationInMinutes": 720,
                        "Challenges": [
                            "OTP,SMS,EMAIL,OATH"
                        ],
                        "AdditionalData": {}
                    },
                    {
                        "Uuid": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                        "Name": "LAB 2FA Authentication Profile",
                        "DurationInMinutes": 0,
                        "Challenges": [
                            "UP",
                            "OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
                        ],
                        "AdditionalData": {
                            "NumberOfQuestions": 1
                        }
                    },
                    {
                        "Uuid": "30804754-3b87-4862-a39e-0f042825a3a0",
                        "Name": "LAB Step-up Authentication Profile",
                        "DurationInMinutes": 0,
                        "Challenges": [
                            "OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
                        ],
                        "AdditionalData": {
                            "NumberOfQuestions": 1
                        }
                    }
                ],
                "PolicyModifiers": [
                    "__centrify_cagent"
                ],
                "DirectoryServices": [
                    {
                        "DisplayNameShort": "Centrify Directory",
                        "directoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E"
                    },
                ],
                "Description": "",
                "Path": "/Policy/LAB User Login Policy",
                "RiskAnalysisLevels": {}
            },
            "Message": null,
            "MessageID": null,
            "Exception": null,
            "ErrorID": null,
            "ErrorCode": null,
            "IsSoftError": false,
            "InnerExceptions": null
        }

	Create policy
	https://developer.centrify.com/reference#post_policy-savepolicyblock3

		Request body format
        {
            "plinks": [
                {
                    "Params": [
                        "d958fad8_f90a_4c40_b986_f6fa31713bba",
                        "d06fc516_8c9b_4f76_a08d_797ca6fd0a55"
                    ],
                    "ID": "/Policy/LAB Machine Login Policy",
                    "EnableCompliant": true,
                    "Description": "",
                    "LinkType": "Role",
                    "PolicySet": "/Policy/LAB Machine Login Policy"
                },
                {
                    "Params": [
                        "sysadmin"
                    ],
                    "ID": "/Policy/LAB System Administrator Login Policy",
                    "EnableCompliant": true,
                    "Description": "Login policy for default PAS system administrator. Only password authentication is required and no 2FA for easy usage.",
                    "LinkType": "Role",
                    "PolicySet": "/Policy/LAB System Administrator Login Policy"
                },
                {
                    "Params": [],
                    "ID": "/Policy/Default Policy",
                    "EnableCompliant": true,
                    "I18NDescriptionTag": "_I18N_DefaultGlobalPolicyDescriptionTag",
                    "Description": "Default Policy Settings.",
                    "LinkType": "Inactive",
                    "PolicySet": "/Policy/Default Policy"
                }
            ],
            "policy": {
                "Path": "/Policy/PolicySet_1",
                "Version": 1,
                "Description": "PolicySet_1",
                "Settings": {
                    "AuthenticationEnabled": true,
                    "/Core/Authentication/AuthenticationRulesDefaultProfileId": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/Authentication/CookieSessionLifespanHours": 24,
                    "/Core/Authentication/CookieAllowPersist": true,
                    "/Core/Authentication/CookiePersistDefault": true,
                    "/Core/Authentication/CookiePersistLifespanHours": 6,
                    "/Core/Authentication/AllowIwa": true,
                    "/Core/Authentication/IwaSetKnownEndpoint": true,
                    "/Core/Authentication/IwaSatisfiesAllMechs": true,
                    "/Core/Authentication/AllowZso": true,
                    "/Core/Authentication/ZsoSkipChallenge": true,
                    "/Core/Authentication/ZsoSetKnownEndpoint": true,
                    "/Core/Authentication/ZsoSatisfiesAllMechs": true,
                    "/Core/Authentication/NoMfaMechLogin": true,
                    "/Core/Authentication/FederatedLoginAllowsMfa": true,
                    "/Core/MfaRestrictions/BlockMobileMechsOnMobileLogin": false,
                    "/Core/Authentication/ContinueFailedSessions": true,
                    "/Core/Authentication/SkipMechsInFalseAdvance": true,
                    "/Core/Authentication/AllowLoginMfaCache": true,
                    "/Core/__centrify_cagent/AuthenticationEnabled": true,
                    "/Core/__centrify_cagent/Authentication/AuthenticationRulesDefaultProfileId": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/__centrify_cagent/Authentication/NoMfaMechLogin": true,
                    "/Core/Css/AuthenticationEnabled": true,
                    "/Core/Css/MfaLogin/AuthenticationRulesDefaultProfileId": "08ccf738-9277-479e-83f4-374609dd9f0e",
                    "/Core/Css/MfaLogin/CssPinningMode": 1,
                    "/Core/Css/WindowsEndpointAuthenticationEnabled": true,
                    "/Core/Css/WinClient/AuthenticationRulesDefaultProfileId": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/Css/PrivilegeElevationEnabled": true,
                    "/Core/Css/Dzdo/AuthenticationRulesDefaultProfileId": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "PasswordResetEnabled": true,
                    "/Core/PasswordReset/PasswordResetEnabled": true,
                    "/Core/PasswordReset/PasswordResetADEnabled": true,
                    "/Core/PasswordReset/PasswordResetIdentityCookieOnly": true,
                    "/Core/PasswordReset/PasswordResetRequiresMfaRestart": true,
                    "/Core/PasswordReset/PasswordResetAuthProfile": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                    "/Core/PasswordReset/PasswordResetMaxAttemptsPerSession": 9,
                    "/Core/PasswordReset/AccountUnlockEnabled": true,
                    "/Core/PasswordReset/AccountUnlockADEnabled": true,
                    "/Core/PasswordReset/AccountUnlockIdentityCookieOnly": true,
                    "/Mobile/EndpointAgent/showAccountLocked": true,
                    "/Core/PasswordReset/AccountUnlockAuthProfile": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/PasswordReset/UseADAdmin": false,
                    "/Core/PasswordReset/Max": 5,
                    "/Core/PasswordReset/MaxTime": 80,
                    "/Core/Security/CDS/PasswordPolicy/MinLength": 9,
                    "/Core/Security/CDS/PasswordPolicy/MaxLength": 22,
                    "/Core/Security/CDS/PasswordPolicy/RequireDigit": true,
                    "/Core/Security/CDS/PasswordPolicy/RequireMixCase": true,
                    "/Core/Security/CDS/PasswordPolicy/RequireSymbol": true,
                    "/Core/Security/CDS/PasswordPolicy/ShowPasswordComplexity": true,
                    "/Core/Security/CDS/PasswordPolicy/NonCdsComplexityHint": "Whatever requirements",
                    "/Core/Security/CDS/PasswordPolicy/AllowRepeatedChar": 3,
                    "/Core/Security/CDS/PasswordPolicy/CheckWeakPassword": true,
                    "/Core/Security/CDS/PasswordPolicy/AllowIncludeUsername": true,
                    "/Core/Security/CDS/PasswordPolicy/AllowIncludeDisplayname": true,
                    "/Core/Security/CDS/PasswordPolicy/RequireUnicode": true,
                    "/Core/Security/CDS/PasswordPolicy/MinAgeInDays": 5,
                    "/Core/Security/CDS/PasswordPolicy/AgeInDays": 445,
                    "/Core/Security/CDS/PasswordPolicy/History": 9,
                    "/Core/PasswordReset/NotifySoft": 49,
                    "/Core/PasswordReset/NotifyHard": 120,
                    "/Core/PasswordChange/NotifyOnMobile": true,
                    "/Core/Security/CDS/LockoutPolicy/Threshold": 9,
                    "/Core/Security/CDS/LockoutPolicy/Window": 54,
                    "/Core/Security/CDS/LockoutPolicy/Duration": 56,
                    "/Core/Security/CDS/ExternalMFA/ShowQRCode": true,
                    "/Core/Authentication/AllowRadius": true,
                    "/Core/Authentication/RadiusUseChallenges": true,
                    "/Core/Authentication/RadiusChallengeProfile": "624a5e94-da4c-4ba3-a596-f56b60637976",
                    "/Core/Authentication/SendRadiusVendorSpecificAttributes": false,
                    "/Core/Authentication/AllowExternalRadius": true,
                    "/Core/PasswordChange/UserChangeAllow": true,
                    "/Core/Authentication/UserUpdateProfile/Password": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/Security/CDS/ExternalMFA/ShowU2f": true,
                    "/Core/Security/CDS/ExternalMFA/U2fUiPrompt": "FIDO2 Security Key",
                    "/Core/Authentication/UserUpdateProfile/U2F": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                    "/Core/Security/CDS/ExternalMFA/ShowQRCodeForSelfService": true,
                    "/Core/Security/CDS/ExternalMFA/UiPrompt": "OATH OTP Client",
                    "/Core/Authentication/UserUpdateProfile/OathProfile": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                    "/Core/Authentication/ConfigureSecurityQuestions": true,
                    "/Core/Authentication/SecurityQuestionPreventDupAnswers": false,
                    "/Core/Authentication/UserSecurityQuestionsPerUser": 2,
                    "/Core/Authentication/AdminSecurityQuestionsPerUser": 3,
                    "/Core/Authentication/SecurityQuestionAnswerMinLength": 4,
                    "/Core/Authentication/UserUpdateProfile/SecurityQuestion": "30804754-3b87-4862-a39e-0f042825a3a0",
                    "/Core/PhoneAuth/UserChangeAllow": true,
                    "/Core/Authentication/MinPhonePinLength": "7",
                    "/Core/Authentication/UserUpdateProfile/PhonePin": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                    "/Core/Security/CDS/AllowUserChangeMFARedirect": true,
                    "/Core/Authentication/UserUpdateProfile/Profile": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb",
                    "/Core/Policy/Culture": "en",
                    "/Mobile/EnrollRules/Common/AllowEnrollment": true,
                    "/Mobile/EnrollRules/Common/AllowJailBrokenDevices": true,
                    "/Mobile/DeviceManagement/EnableInviteBasedEnrollment": true,
                    "/Mobile/Software/Policies/Centrify/Common/AllowNotificationOnMutipleDevices": true,
                    "/Mobile/Software/Policies/Centrify/iOSSettings/AllowDebugLogging": true,
                    "/Mobile/Software/Policies/Centrify/iOSSettings/Restrictions/LocationTracking": true,
                    "/Mobile/Software/Policies/Centrify/Application/Security/MobileAuthenticator/ForceFingerprintForMobileAuthenticator": true,
                    "/Mobile/Software/Policies/Centrify/Application/Security/MobileAuthenticator/ForceFingerprintForMobileAuthenticatorAllowFallbackAppPin": true,
                    "/Mobile/Software/Policies/Centrify/Application/Passcode/ForceAppPin": true,
                    "/Mobile/Software/Policies/Centrify/Application/Passcode/AppInactivityTimeout": 15,
                    "/Mobile/Software/Policies/Centrify/Application/Passcode/AppLockOnExit": true
                },
                "Newpolicy": true
            }
        }

		Respond result
        {
            "success": true,
            "Result": {
                "RevStamp": "637331588600000000"
            },
            "Message": null,
            "MessageID": null,
            "Exception": null,
            "ErrorID": null,
            "ErrorCode": null,
            "IsSoftError": false,
            "InnerExceptions": null
        }


	Update policy
	https://developer.centrify.com/reference#post_roles-updaterole-1

		Request body format
		{
			"policy": {
				"Version": 0,
				"Settings": {
					"/Core/__centrify_cagent/Authentication/AuthenticationRulesHighAuthRequestedProfileId": "AlwaysAllowed",
					"/Core/__centrify_cagent/AuthenticationEnabled": true,
					"/Core/__centrify_cagent/Authentication/ChallengeDefinitionId": "7f7fcf81-ee3a-4b78-8e23-97da87f5ac7c",
					"/Core/__centrify_cagent/Authentication/AuthenticationRules": {
						"_UniqueKey": "Condition",
						"_Value": [],
						"Enabled": true,
						"_Type": "RowSet"
					},
					"/Core/Css/MfaLogin/CssPinningMode": 0,
					"/Core/__centrify_cagent/Authentication/AuthenticationRulesDefaultProfileId": "30804754-3b87-4862-a39e-0f042825a3a0"
				},
				"RevStamp": "637334975380000000",
				"RadiusClientList": [],
				"AuthProfiles": [
					{
						"Uuid": "cf484a69-4f1b-42e0-9286-2e0ce972b362",
						"Name": "Default New Device Login Profile",
						"DurationInMinutes": 720,
						"Challenges": [
							"UP",
							"OTP,SMS,EMAIL,OATH,SQ"
						],
						"AdditionalData": {
							"NumberOfQuestions": 1
						}
					}
				],
				"PolicyModifiers": [
					"__centrify_cagent"
				],
				"DirectoryServices": [
					{
						"DisplayNameShort": "Centrify Directory",
						"directoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E"
					}
				],
				"Description": "Test Policy 1",
				"Path": "/Policy/Test Policy 1",
				"RiskAnalysisLevels": {},
				"Newpolicy": false
			},
			"plinks": [
				{
					"Params": [
						"d958fad8_f90a_4c40_b986_f6fa31713bba",
						"d06fc516_8c9b_4f76_a08d_797ca6fd0a55"
					],
					"ID": "/Policy/LAB Machine Login Policy",
					"EnableCompliant": true,
					"Description": "",
					"LinkType": "Role",
					"PolicySet": "/Policy/LAB Machine Login Policy"
				},
				{
					"Params": [],
					"ID": "/Policy/Default Policy",
					"EnableCompliant": true,
					"I18NDescriptionTag": "_I18N_DefaultGlobalPolicyDescriptionTag",
					"Description": "Default Policy Settings.",
					"LinkType": "Inactive",
					"PolicySet": "/Policy/Default Policy"
				}
			]
		}

		Respond result


	Delete policy
	https://developer.centrify.com/reference#post_policy-deletepolicyblock

		Request body format
        {
            "path": "/Policy/PolicySet_1"
        }

		Respond result
        {
            "success": true,
            "Result": null,
            "Message": null,
            "MessageID": null,
            "Exception": null,
            "ErrorID": null,
            "ErrorCode": null,
            "IsSoftError": false,
            "InnerExceptions": null
        }

*/