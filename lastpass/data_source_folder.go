package lastpass

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rezroo/terraform-provider-lastpass/api"
)

// DataSourceFolder describes our lastpass folder data source
func DataSourceFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceFolderRead,
		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			"folder": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// DataSourceFolderRead reads resource from upstream/lastpass
func DataSourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Starting folder read")

	var folder_name = d.Get("folder").(string)

	// log.Printf("[INFO] before read")
	data, err := dataSourceFolderRead(m, folder_name)
	// log.Printf("[INFO] after read")

	if err != nil {
		return diag.FromErr(err)
	} else if len(data) == 0 {
		d.SetId("")
		return diags
	}

	d.SetId(folder_name)

	return diags
}

func dataSourceFolderRead(m interface{}, folder string) ([]api.FolderShare, error) {
	client := m.(*api.Client)
	data, err := client.ReadShares(folder)

	log.Printf("[INFO] client.ReadShares complete with", len(data))

	if err != nil {
		log.Printf("[ERROR]", err)
		var err = errors.New("Error in dataSourceFolderRead")
		return data, err
	}

	// else if len(data) > 1 {
	//     return []api.FolderShare{}, err
	// }

	log.Printf("[INFO] dataSourceFolderRead finished")
	return data, nil
}
