package genesyscloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRoutingQueueBasic(t *testing.T) {
	var (
		queueResource = "test-queue"
		queueName     = "Terraform Test Queue-" + uuid.NewString()
		queueDesc     = "This is a test"

		queueDataSource = "test-queue-ds"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				// Create
				Config: GenerateRoutingQueueResource(
					queueResource,
					queueName,
					queueDesc,
					nullValue, // MANDATORY_TIMEOUT
					"200000",  // acw_timeout
					nullValue, // ALL
					nullValue, // auto_answer_only true
					nullValue, // No calling party name
					nullValue, // No calling party number
					nullValue, // enable_manual_assignment false
					nullValue, // enable_transcription false
				) + generateRoutingQueueDataSource(
					queueDataSource,
					"genesyscloud_routing_queue."+queueResource+".name",
					"genesyscloud_routing_queue."+queueResource,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.genesyscloud_routing_queue."+queueDataSource,
						"id", "genesyscloud_routing_queue."+queueResource, "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceRoutingQueueCaching(t *testing.T) {
	var (
		queue1ResourceId = "queue1"
		queueName1       = "terraform test queue " + uuid.NewString()
		queue2ResourceId = "queue2"
		queueName2       = "terraform test queue " + uuid.NewString()
		queue3ResourceId = "queue3"
		queueName3       = "terraform test queue " + uuid.NewString()

		dataSource1Id = "data-1"
		dataSource2Id = "data-2"
		dataSource3Id = "data-3"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				Config: generateRoutingQueueResourceBasic( // queue resource
					queue1ResourceId,
					queueName1,
				) + generateRoutingQueueResourceBasic( // queue resource
					queue2ResourceId,
					queueName2,
				) + generateRoutingQueueResourceBasic( // queue resource
					queue3ResourceId,
					queueName3,
				) + generateRoutingQueueDataSource( // queue data source
					dataSource1Id,
					strconv.Quote(queueName1),
					"genesyscloud_routing_queue."+queue1ResourceId,
				) + generateRoutingQueueDataSource( // queue data source
					dataSource2Id,
					strconv.Quote(queueName2),
					"genesyscloud_routing_queue."+queue2ResourceId,
				) + generateRoutingQueueDataSource( // queue data source
					dataSource3Id,
					strconv.Quote(queueName3),
					"genesyscloud_routing_queue."+queue3ResourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("genesyscloud_routing_queue."+queue1ResourceId, "id",
						"data.genesyscloud_routing_queue."+dataSource1Id, "id"),
					resource.TestCheckResourceAttrPair("genesyscloud_routing_queue."+queue2ResourceId, "id",
						"data.genesyscloud_routing_queue."+dataSource2Id, "id"),
					resource.TestCheckResourceAttrPair("genesyscloud_routing_queue."+queue3ResourceId, "id",
						"data.genesyscloud_routing_queue."+dataSource3Id, "id"),
				),
			},
		},
		CheckDestroy: testVerifyQueuesDestroyed,
	})
}

func generateRoutingQueueDataSource(
	resourceID string,
	name string,
	// Must explicitly use depends_on in terraform v0.13 when a data source references a resource
	// Fixed in v0.14 https://github.com/hashicorp/terraform/pull/26284
	dependsOnResource string) string {
	return fmt.Sprintf(`data "genesyscloud_routing_queue" "%s" {
		name = %s
		depends_on=[%s]
	}
	`, resourceID, name, dependsOnResource)
}
