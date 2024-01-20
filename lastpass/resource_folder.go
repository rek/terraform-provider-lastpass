package lastpass

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rezroo/terraform-provider-lastpass/api"
)

// ResourceFolder describes lastpass_Folder resource
// Reference: https://github.com/lastpass/lastpass-cli/blob/master/cmd-share.c
func ResourceSharedFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext:   ResourceSharedFolderRead,
		CreateContext: ResourceSharedFolderCreate,
		UpdateContext: ResourceSharedFolderUpdate,
		DeleteContext: ResourceSharedFolderDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceSharedFolderImporter,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"folder": {
				Type:     schema.TypeString,
				Required: true,
				// Optional: true,
				// ForceNew: true,
				// Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hide": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// "accept": {
			// 	Type:     schema.TypeBool,
			// 	Optional:    true,
			// },
		},
	}
}

// ResourceSharedFolderRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceSharedFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Starting ResourceSharedFolderRead")
	var diags diag.Diagnostics

	client := m.(*api.Client)
	data, err := client.ReadShare(d.Id())
	// data, err := dataSourceFolderRead(m, d.Id())

	log.Printf("[INFO] ResourceSharedFolderRead response", data)

	if err != nil {
		return diag.FromErr(err)
	}
	// if len(data) == 0 {
	// 	d.SetId("")
	// 	return nil
	// }

	// if len(data) > 2 {
	// 	var err = errors.New("Too many datas!")
	// 	return diag.FromErr(err)
	// }

	folder_name, _ := splitId(d.Id())

	d.Set("id", d.Id())
	d.Set("folder", folder_name)
	d.Set("user", data.Name)
	d.Set("email", data.Email)
	d.Set("read_only", data.ReadOnly)
	d.Set("admin", data.Admin)
	d.Set("hide", data.Hide)

	return diags
}

func ResourceSharedFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics
	err := client.DeleteFolder(d.Id(), d.Get("email").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func ResourceSharedFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	folder_share := api.FolderShare{
		Id:       d.Id(),
		Folder:   d.Get("folder").(string),
		Name:     d.Get("user").(string),
		Email:    d.Get("email").(string),
		ReadOnly: d.Get("read_only").(bool),
		Admin:    d.Get("admin").(bool),
		Hide:     d.Get("hide").(bool),
	}
	client := m.(*api.Client)
	err := client.UpdateFolder(folder_share)
	if err != nil {
		return diag.FromErr(err)
	}
	return DataSourceFolderShareRead(ctx, d, m)
}

func ResourceSharedFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics

	folder_share := api.FolderShare{
		Id:       d.Id(),
		Folder:   d.Get("folder").(string),
		Name:     d.Get("user").(string),
		Email:    d.Get("email").(string),
		ReadOnly: d.Get("read_only").(bool),
		Admin:    d.Get("admin").(bool),
		Hide:     d.Get("hide").(bool),
	}

	err := client.CreateFolder(folder_share)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// called to import an existing resource.
func ResourceSharedFolderImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] Starting import")

	// take 2
	folder_name, _ := splitId(d.Id())

	// // test that both of these exist
	// if folder_name == "" || email == "" {
	// 	var error = errors.New("Invalid ID. Use the format: folder/email")
	// 	diag.FromErr(error)
	// 	return nil, error
	// }
	// client := m.(*api.Client)
	// data, err := client.ReadShares(folder_name)

	// take 1
	// data, err := dataSourceFolderRead(m, d.Id())

	client := m.(*api.Client)
	user, err := client.ReadShare(d.Id())

	if err != nil {
		return nil, err
	}
	if user.Email == "" {
		// if len(data) == 0 {
		var err = errors.New("No shared data found for this folder ID")
		return nil, err
	}
	// log.Printf("[INFO] Read success, total: %v", len(data))

	// var user = findUser(data, email)

	if user.Email == "" {
		var err = errors.New("Failed to find record for '%s'" + d.Id())
		diag.FromErr(err)
		return nil, err
	}

	log.Printf("[INFO] --matched==", user.Email)

	d.Set("id", d.Id())
	d.Set("folder", folder_name)
	d.Set("user", user.Name)
	d.Set("email", user.Email)
	d.Set("read_only", user.ReadOnly)
	d.Set("admin", user.Admin)
	d.Set("hide", user.Hide)

	return []*schema.ResourceData{d}, nil
}

// id looks like: folder/email
func splitId(id string) (string, string) {
	splitLine := strings.Split(id, "/")
	var folder_name = splitLine[0]
	var email = splitLine[1]
	return folder_name, email
}

func findUser(shares []api.FolderShare, email string) api.FolderShare {
	share_map := make(map[string]api.FolderShare)
	for _, share := range shares {
		share_map[share.Email] = share
	}

	// _, found := share_map[email]

	return share_map[email]
}
