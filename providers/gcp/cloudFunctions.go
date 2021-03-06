// Copyright 2018 The Terraformer Authors.
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

package gcp

import (
	"context"
	"log"
	"strings"

	"google.golang.org/api/cloudfunctions/v1"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"golang.org/x/oauth2/google"
)

var cloudFunctionsAllowEmptyValues = []string{""}

var cloudFunctionsAdditionalFields = map[string]string{}

type CloudFunctionsGenerator struct {
	GCPService
}

// Run on CloudFunctionsList and create for each TerraformResource
func (g CloudFunctionsGenerator) createResources(functionsList *cloudfunctions.ProjectsLocationsFunctionsListCall, ctx context.Context) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}
	if err := functionsList.Pages(ctx, func(page *cloudfunctions.ListFunctionsResponse) error {
		for _, functions := range page.Functions {
			t := strings.Split(functions.Name, "/")
			name := t[len(t)-1]
			resources = append(resources, terraform_utils.NewResource(
				g.GetArgs()["project"]+"/"+g.GetArgs()["region"]+"/"+name,
				g.GetArgs()["region"]+"_"+name,
				"google_cloudfunctions_function",
				"google",
				map[string]string{},
				cloudFunctionsAllowEmptyValues,
				cloudFunctionsAdditionalFields,
			))
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return resources
}

// Generate TerraformResources from GCP API,
// from each CloudFunctions create 1 TerraformResource
// Need CloudFunctions name as ID for terraform resource
func (g *CloudFunctionsGenerator) InitResources() error {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudfunctions.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudfunctionsService, err := cloudfunctions.New(c)
	if err != nil {
		log.Fatal(err)
	}

	functionsList := cloudfunctionsService.Projects.Locations.Functions.List("projects/" + g.GetArgs()["project"] + "/locations/" + g.GetArgs()["region"])

	g.Resources = g.createResources(functionsList, ctx)
	g.PopulateIgnoreKeys()
	return nil

}
