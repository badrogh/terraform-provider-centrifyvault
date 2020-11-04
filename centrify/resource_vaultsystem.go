package centrify

import (
	"fmt"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceVaultSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultSystemCreate,
		Read:   resourceVaultSystemRead,
		Update: resourceVaultSystemUpdate,
		Delete: resourceVaultSystemDelete,
		Exists: resourceVaultSystemExists,

		Schema: map[string]*schema.Schema{
			// System -> Settings menu related settings
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the system",
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Hostname or IP address of the system",
				ValidateFunc: validation.NoZeroValues,
			},
			"computer_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the system",
				ValidateFunc: validation.StringInSlice([]string{
					"Windows",
					"Unix",
					"CiscoAsyncOS",
					"CiscoIOS",
					"CiscoNXOS",
					"JuniperJunos",
					"HpNonStopOS",
					"IBMi",
					"CheckPointGaia",
					"PaloAltoNetworksPANOS",
					"F5NetworksBIGIP",
					"VMwareVMkernel",
					"GenericSsh",
					"CustomSsh",
				}, false),
			},
			"session_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login session type that the system supports",
				ValidateFunc: validation.StringInSlice([]string{
					"Rdp",
					"Ssh",
				}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the system",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Port that used to connect to the system",
				ValidateFunc: validation.IsPortNumber,
			},
			"use_my_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable Use My Account",
			},
			"management_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Management mode of the system. For Windows only",
				ValidateFunc: validation.StringInSlice([]string{
					"RpcOverTcp",
					"Smb",
					"WinRMOverHttp",
					"WinRMOverHttps",
					"Disabled",
				}, false),
			},
			"management_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Management port for account management. For Windows only",
				ValidateFunc: validation.IsPortNumber,
			},
			"system_timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true, // Set it to computed as once it is set it can't be unset. It causes TF always think there is change
				Description: "System time zone",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxyuser": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxyuser_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"proxyuser_managed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// System -> Policy menu related settings
			"checkout_lifetime": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Specifies the number of minutes that a checked out password is valid.",
				ValidateFunc: validation.IntAtLeast(15),
			},
			"allow_remote_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow access from a public network (web client only)",
			},
			"allow_rdp_clipboard": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow RDP client to sync local clipboard with remote session",
			},
			"default_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default System Login Profile (used if no conditions matched)",
			},
			// System -> Advanced menu related settings
			"local_account_automatic_maintenance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Local Account Automatic Maintenance",
			},
			"local_account_manual_unlock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Local Account Manual Unlock",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AD domain that this system belongs to",
			},
			"allow_multiple_checkouts": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow multiple password checkouts for this system",
			},
			"enable_password_rotation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable periodic password rotation",
			},
			"password_rotate_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Password rotation interval (days)",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"enable_password_rotation_after_checkin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable password rotation after checkin",
			},
			"minimum_password_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum Password Age (days)",
			},
			"password_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
				//Computed:    true, // we want to remove this setting if it is not set so do not set to computed
				Description: "Password complexity profile id",
			},
			"enable_password_history_cleanup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable periodic password history cleanup",
			},
			"password_historycleanup_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Password history cleanup (days)",
				ValidateFunc: validation.IntAtLeast(90),
			},
			"enable_sshkey_rotation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable periodic SSH key rotation",
			},
			"sshkey_rotate_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "SSH key rotation interval (days)",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"minimum_sshkey_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum SSH Key Age (days)",
			},
			"sshkey_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSH Key Generation Algorithm",
				ValidateFunc: validation.StringInSlice([]string{
					"RSA_1024",
					"RSA_2048",
					"ECDSA_P256",
					"ECDSA_P384",
					"ECDSA_P521",
					"EdDSA_Ed448",
					"EdDSA_Ed25519",
				}, false),
			},
			"enable_sshkey_history_cleanup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable periodic SSH key history cleanup",
			},
			"sshkey_historycleanup_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "SSH key history cleanup (days)",
				ValidateFunc: validation.IntAtLeast(90),
			},

			// System -> Zone Role Workflow menu related settings
			"use_domainadmin_for_zonerole_workflow": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use Domain Administrator Account for Zone Role Workflow operations",
			},
			"enable_zonerole_workflow": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable zone role requests for this system",
			},
			// System -> Connectors menu related settings
			"connector_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of Connectors",
			},
			// Add to Sets
			"sets": {
				Type:     schema.TypeSet,
				Optional: true,
				//Computed: true,
				Set: schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Add to list of Sets",
			},
			"permission":     getPermissionSchema(),
			"challenge_rule": getChallengeRulesSchema(),
			/* Can't use TypeMap because some attributes are string but some are bool
			"proxy_account": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			*/
			//"proxy_account": proxyAccountSchema(),
		},
	}
}

func resourceVaultSystemExists(d *schema.ResourceData, m interface{}) (bool, error) {
	LogD.Printf("Checking System exist: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewVaultSystem(client)
	object.ID = d.Id()
	err := object.Read()

	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	LogD.Printf("System exists in tenant: %s", object.ID)
	return true, nil
}

func resourceVaultSystemRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Reading System: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	// Create a System object and populate ID attribute
	object := NewVaultSystem(client)
	object.ID = d.Id()
	err := object.Read()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading System: %v", err)
	}
	//LogD.Printf("System from tenant: %v", object)

	schemamap, err := generateSchemaMap(object)
	if err != nil {
		return err
	}
	LogD.Printf("Generated Map for resourceSystemRead(): %+v", schemamap)
	for k, v := range schemamap {
		if k == "connector_list" {
			// Convert "value1,value1" to schema.TypeSet
			d.Set("connector_list", schema.NewSet(schema.HashString, StringSliceToInterface(strings.Split(v.(string), ","))))
		} else {
			d.Set(k, v)
		}
	}

	LogD.Printf("Completed reading System: %s", object.Name)
	return nil
}

func resourceVaultSystemCreate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning System creation: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)

	// Create a System object and populate all attributes
	object := NewVaultSystem(client)
	err := createUpateGetSystemData(d, object)
	if err != nil {
		return err
	}

	resp, err := object.Create()
	if err != nil {
		return fmt.Errorf("Error creating System: %v", err)
	}

	id := resp.Result
	if id == "" {
		return fmt.Errorf("System ID is not set")
	}
	d.SetId(id)
	// Need to populate ID attribute for subsequence processes
	object.ID = id

	d.SetPartial("name")
	d.SetPartial("fqdn")
	d.SetPartial("computer_class")
	d.SetPartial("session_type")
	d.SetPartial("description")

	// 2nd step to update system login profile
	// Create API call doesn't set system login profile so need to run update again
	if object.LoginDefaultProfile != "" {
		LogD.Printf("Update login profile for System creation: %s", ResourceIDString(d))
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating System attribute: %v", err)
		}
		d.SetPartial("default_profile_id")
	}

	// 3rd step to add system to Sets
	if len(object.Sets) > 0 {
		for _, v := range object.Sets {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Server"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "add")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error adding System to Set: %v", err)
			}
		}
		d.SetPartial("sets")
	}

	// 4th step to add permissions
	if _, ok := d.GetOk("permission"); ok {
		_, err = object.SetPermissions(false)
		if err != nil {
			return fmt.Errorf("Error setting System permissions: %v", err)
		}
		d.SetPartial("permission")
	}

	// Creation completed
	d.Partial(false)
	LogD.Printf("Creation of System completed: %s", object.Name)
	return resourceVaultSystemRead(d, m)
}

func resourceVaultSystemUpdate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning System update: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)
	object := NewVaultSystem(client)

	object.ID = d.Id()
	err := createUpateGetSystemData(d, object)
	if err != nil {
		return err
	}

	// Deal with normal attribute changes first
	if d.HasChanges("name", "fqdn", "description", "port", "use_my_account", "management_mode", "system_timezone", "proxyuser", "proxyuser_password",
		"proxyuser_managed", "checkout_lifetime", "allow_remote_access", "allow_rdp_clipboard", "default_profile_id",
		"local_account_automatic_maintenance", "local_account_manual_unlock", "domain_id", "allow_multiple_checkouts", "enable_password_rotation",
		"password_rotate_interval", "enable_password_rotation_after_checkin", "minimum_password_age", "password_profile_id", "enable_password_history_cleanup",
		"password_historycleanup_duration", "enable_sshkey_rotation", "sshkey_rotate_interval", "minimum_sshkey_age", "sshkey_algorithm",
		"enable_sshkey_history_cleanup", "sshkey_historycleanup_duration", "use_domainadmin_for_zonerole_workflow", "enable_zonerole_workflow",
		"choose_connector", "connector_list", "challenge_rule") {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating System attribute: %v", err)
		}
		LogD.Printf("Updated attributes to: %+v", object)
		d.SetPartial("name")
		d.SetPartial("fqdn")
		d.SetPartial("computer_class")
		d.SetPartial("session_type")
		d.SetPartial("description")
		d.SetPartial("default_profile_id")
		d.SetPartial("challenge_rule")
	}

	// Deal with Set member
	if d.HasChange("sets") {
		old, new := d.GetChange("sets")
		// Remove old Sets
		for _, v := range flattenSchemaSetToStringSlice(old) {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Server"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "remove")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error removing System from Set: %v", err)
			}
		}
		// Add new Sets
		for _, v := range flattenSchemaSetToStringSlice(new) {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Server"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "add")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error adding System to Set: %v", err)
			}
		}
		d.SetPartial("sets")
	}

	// Deal with Permissions
	if d.HasChange("permission") {
		old, new := d.GetChange("permission")
		// We don't want to care the details of changes
		// So, let's first remove the old permissions
		var err error
		var perms map[string]string
		if object.ComputerClass == "Windows" || object.ComputerClass == "Unix" {
			perms = winnixPermissions
		} else {
			perms = systemPermissions
		}
		if old != nil {
			// do not validate old values
			object.Permissions, err = expandPermissions(old, perms, false)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(true)
			if err != nil {
				return fmt.Errorf("Error removing System permissions: %v", err)
			}
		}

		if new != nil {
			object.Permissions, err = expandPermissions(new, perms, true)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(false)
			if err != nil {
				return fmt.Errorf("Error adding System permissions: %v", err)
			}
		}
		d.SetPartial("permission")
	}

	d.Partial(false)
	LogD.Printf("Updating of System completed: %s", object.Name)
	return resourceVaultSystemRead(d, m)
}

func resourceVaultSystemDelete(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning deletion of System: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewVaultSystem(client)
	object.ID = d.Id()
	resp, err := object.Delete()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		return fmt.Errorf("Error deleting System: %v", err)
	}

	if resp.Success {
		d.SetId("")
	}

	LogD.Printf("Deletion of System completed: %s", ResourceIDString(d))
	return nil
}

func createUpateGetSystemData(d *schema.ResourceData, object *VaultSystem) error {
	// System -> Settings menu related settings
	object.Name = d.Get("name").(string)
	object.FQDN = d.Get("fqdn").(string)
	object.ComputerClass = d.Get("computer_class").(string)
	object.SessionType = d.Get("session_type").(string)
	if v, ok := d.GetOk("description"); ok {
		object.Description = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		object.Port = v.(int)
	}
	if v, ok := d.GetOk("use_my_account"); ok {
		object.UseMyAccount = v.(bool)
	}
	if v, ok := d.GetOk("management_mode"); ok {
		object.ManagementMode = v.(string)
	}
	if v, ok := d.GetOk("system_timezone"); ok {
		object.TimeZoneID = v.(string)
	}
	if v, ok := d.GetOk("proxyuser"); ok {
		object.ProxyUser = v.(string)
	}
	if v, ok := d.GetOk("proxyuser_password"); ok {
		object.ProxyUserPassword = v.(string)
	}
	if v, ok := d.GetOk("proxyuser_managed"); ok {
		object.ProxyUserIsManaged = v.(bool)
	}
	// System -> Policy menu related settings
	if v, ok := d.GetOk("checkout_lifetime"); ok {
		object.DefaultCheckoutTime = v.(int)
	}
	if v, ok := d.GetOk("allow_remote_access"); ok {
		object.AllowRemote = v.(bool)
	}
	if v, ok := d.GetOk("allow_rdp_clipboard"); ok {
		object.AllowRdpClipboard = v.(bool)
	}
	if v, ok := d.GetOk("default_profile_id"); ok {
		object.LoginDefaultProfile = v.(string)
	}
	// System -> Advanced menu related settings
	if v, ok := d.GetOk("local_account_automatic_maintenance"); ok {
		object.AllowAutomaticLocalAccountMaintenance = v.(bool)
	}
	if v, ok := d.GetOk("local_account_manual_unlock"); ok {
		object.AllowManualLocalAccountUnlock = v.(bool)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		object.DomainID = v.(string)
	}
	if v, ok := d.GetOk("allow_multiple_checkouts"); ok {
		object.AllowMultipleCheckouts = v.(bool)
	}
	if v, ok := d.GetOk("enable_password_rotation"); ok {
		object.AllowPasswordRotation = v.(bool)
	}
	if v, ok := d.GetOk("password_rotate_interval"); ok {
		object.PasswordRotateDuration = v.(int)
	}
	if v, ok := d.GetOk("enable_password_rotation_after_checkin"); ok {
		object.AllowPasswordRotationAfterCheckin = v.(bool)
	}
	if v, ok := d.GetOk("minimum_password_age"); ok {
		object.MinimumPasswordAge = v.(int)
	}
	if v, ok := d.GetOk("password_profile_id"); ok {
		object.PasswordProfileID = v.(string)
	}
	if v, ok := d.GetOk("enable_password_history_cleanup"); ok {
		object.AllowPasswordHistoryCleanUp = v.(bool)
	}
	if v, ok := d.GetOk("password_historycleanup_duration"); ok {
		object.PasswordHistoryCleanUpDuration = v.(int)
	}
	if v, ok := d.GetOk("enable_sshkey_rotation"); ok {
		object.AllowSSHKeysRotation = v.(bool)
	}
	if v, ok := d.GetOk("sshkey_rotate_interval"); ok {
		object.SSHKeysRotateDuration = v.(int)
	}
	if v, ok := d.GetOk("minimum_sshkey_age"); ok {
		object.MinimumSSHKeysAge = v.(int)
	}
	if v, ok := d.GetOk("sshkey_algorithm"); ok {
		object.SSHKeysGenerationAlgorithm = v.(string)
	}
	if v, ok := d.GetOk("enable_sshkey_history_cleanup"); ok {
		object.AllowSSHKeysCleanUp = v.(bool)
	}
	if v, ok := d.GetOk("sshkey_historycleanup_duration"); ok {
		object.SSHKeysCleanUpDuration = v.(int)
	}
	// System -> Zone Role Workflow menu related settings
	if v, ok := d.GetOk("use_domainadmin_for_zonerole_workflow"); ok {
		object.DomainOperationsEnabled = v.(bool)
	}
	if v, ok := d.GetOk("enable_zonerole_workflow"); ok {
		object.ZoneRoleWorkflowEnabled = v.(bool)
	}
	// System -> Connectors menu related settings
	if v, ok := d.GetOk("connector_list"); ok {
		object.ProxyCollectionList = flattenSchemaSetToString(v.(*schema.Set))
	}
	// Sets
	if v, ok := d.GetOk("sets"); ok {
		object.Sets = flattenSchemaSetToStringSlice(v)
	}
	// Permissions
	if v, ok := d.GetOk("permission"); ok {
		var err error
		var perms map[string]string
		if object.ComputerClass == "Windows" || object.ComputerClass == "Unix" {
			perms = winnixPermissions
		} else {
			perms = systemPermissions
		}
		object.Permissions, err = expandPermissions(v, perms, true)
		if err != nil {
			return err
		}
	}
	// Challenge rules
	if v, ok := d.GetOk("challenge_rule"); ok {
		object.ChallengeRules = expandChallengeRules(v.([]interface{}))
		// Perform validations
		if err := validateChallengeRules(object.ChallengeRules); err != nil {
			return fmt.Errorf("Schema setting error: %s", err)
		}
	}

	// Perform validations
	if err := object.ValidateZoneWorkflow(); err != nil {
		return fmt.Errorf("Schema setting error: %s", err)
	}
	return nil
}
