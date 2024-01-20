package lastpass

import (
	"context"
	"strings"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rezroo/terraform-provider-lastpass/api"
)

// DataSourceFolderShare describes our lastpass folder share data source
func DataSourceFolderShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceFolderRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hide": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// DataSourceFolderRead reads resource from upstream/lastpass
func DataSourceFolderShareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Starting folder read")

	var id = d.Get("id").(string)

	splitLine := strings.Split(id, "/")
	var folder_name = splitLine[0]
	
	// log.Printf("[INFO] before read")
	data, err := dataSourceFolderShareRead(m, id)
	// log.Printf("[INFO] after read")

	if err != nil {
		return diag.FromErr(err)
	} 

	// else if len(data) == 0 {
	// 	d.SetId("")
	// 	return diags
	// }

	d.SetId(id)

	d.Set("folder", folder_name)
	d.Set("user", data.Name)
	d.Set("email", data.Email)
	d.Set("read_only", data.ReadOnly)
	d.Set("admin", data.Admin)
	d.Set("hide", data.Hide)

	return diags
}

func dataSourceFolderShareRead(m interface{}, id string) (api.FolderShare, error) {
	client := m.(*api.Client)
	data, err := client.ReadShare(id)

	log.Printf("[INFO] client.ReadShare complete for %v", id)

	if err != nil {
		log.Printf("[ERROR]", err)
		var err = errors.New("Error in dataSourceFolderRead")
		return data, err
	}

	// else if len(data) > 1 {
	//     return []api.FolderShare{}, err
	// }

	log.Printf("[INFO] dataSourceFolderShareRead finished")
	return data, nil
}
