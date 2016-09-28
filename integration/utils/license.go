package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func AddValidLicense(t *testing.T, serverURL, username, password string) {
	log.Infof("Attempting to add license on %s", serverURL)
	orcaURL, err := neturl.Parse(serverURL)
	require.Nil(t, err)
	dclient, err := GetUserDockerClient(serverURL, username, password)
	require.Nil(t, err)
	client := dclient.HTTPClient

	orcaURL.Path = "/api/config/license"
	// WARNING! this is real license, don't let it leak into the field
	data := []byte(
		`{"auto_refresh":true,"license_config":{"key_id":"4Hg5DGMH78wN5ZjNjbau_agErRqNE5aQ-R3MnUiNYdGg","private_key":"y2dhoSi4jAYKpclaXnmd1R_RJkYy7ySmKCis9e1JpfH6","authorization":"ewogICAicGF5bG9hZCI6ICJleUpsZUhCcGNtRjBhVzl1SWpvaU1qQTBNeTB3TkMweU4xUXhPRG95T0RvMU1sb2lMQ0owYjJ0bGJpSTZJbUpmZVV0VVRFeFVjVUppTmxCdE5ucE5hRlUxT1ZGRlJUQnFaRVYxUVhoWFIzUmZZM1Z0TVZsb1NqQTlJaXdpYldGNFJXNW5hVzVsY3lJNk1UQXNJbXhwWTJWdWMyVlVlWEJsSWpvaVQyWm1iR2x1WlNJc0luUnBaWElpT2lKRmRtRnNkV0YwYVc5dUluMCIsCiAgICJzaWduYXR1cmVzIjogWwogICAgICB7CiAgICAgICAgICJoZWFkZXIiOiB7CiAgICAgICAgICAgICJqd2siOiB7CiAgICAgICAgICAgICAgICJlIjogIkFRQUIiLAogICAgICAgICAgICAgICAia2V5SUQiOiAiSjdMRDo2N1ZSOkw1SFo6VTdCQToyTzRHOjRBTDM6T0YyTjpKSEdCOkVGVEg6NUNWUTpNRkVPOkFFSVQiLAogICAgICAgICAgICAgICAia2lkIjogIko3TEQ6NjdWUjpMNUhaOlU3QkE6Mk80Rzo0QUwzOk9GMk46SkhHQjpFRlRIOjVDVlE6TUZFTzpBRUlUIiwKICAgICAgICAgICAgICAgImt0eSI6ICJSU0EiLAogICAgICAgICAgICAgICAibiI6ICJ5ZEl5LWxVN283UGNlWS00LXMtQ1E1T0VnQ3lGOEN4SWNRSVd1Szg0cElpWmNpWTY3MzB5Q1lud0xTS1Rsdy1VNlVDX1FSZVdSaW9NTk5FNURzNVRZRVhiR0c2b2xtMnFkV2JCd2NDZy0yVVVIX09jQjlXdVA2Z1JQSHBNRk1zeER6V3d2YXk4SlV1SGdZVUxVcG0xSXYtbXE3bHA1blFfUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PX2dlTFBTLWhnS1B0SUhsa2c2X1djb3hUR29LUDc5ZF93YUhZeEdObDdXaFNuZWlCU3hicGJRQUtrMjFsZzc5OFhiN3ZaeUVBVERNclJSOU1lRTZBZGo1SEpwWTNDb3lSQVBDbWFLR1JDSzR1b1pTb0l1MGhGVmxLVVB5YmJ3MDAwR08td2EyS044VXdnSUltMGk1STF1VzlHa3E0empCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcC1ERFBFOVRHZzR6WW1YakpueFpxSEVkdUdxZGV2WjhYTUkwdWtma0dJSTE0d1VPaU1JSUlyWGxFY0JmXzQ2SThnUVdEenh5Y1plX0pHWC1MQXVheVhyeXJVRmVoVk5VZFpVbDl3WE5hSkIta2FDcXo1UXdhUjkzc0d3LVFTZnREME52TGU3Q3lPSC1FNnZnNlN0X05lVHZndjhZbmhDaVhJbFo4SE9mSXdOZTd0RUZfVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybWs0bFZJWjdWRF9MU1ctaTd5b1N1cnRwc1BYY2UycEtESW8zMGxKR2hPXzNLVW1sMlNVWkNxekoxeUVtS3B5c0g1SERXOWNzSUZDQTNkZUFqZlpVdk43VSIKICAgICAgICAgICAgfSwKICAgICAgICAgICAgImFsZyI6ICJSUzI1NiIKICAgICAgICAgfSwKICAgICAgICAgInNpZ25hdHVyZSI6ICJMWEtUclBfVTJEUGVlWlBZaFlZdjZJTm1BU1dERWYtMEV5ZDdTb1hwdDdVeDRVOU11VVF2dzlCTFotaDQ0R3JYeXVaeGxYMDdjY2xRc0NRNWZLV1JRdy1XQkphQi04UlRqWXFPaXh3UldCZGlFaDM1c0tNSjRpSzFKbkMxLXpPN1JxTkdycmhscGgtZHM5QUhBT3c0THM5REJRWmZFVURzUzl6X296R1liOHlDR0FmTS1EOWsxTFF5djBoREJJd0ZnMDRmSUFiTkZucmlZckRiRzc4WGl6LWdXQlJieExCdXZxV2lnLTFfOGtab1VkMVRpX1JSREowTThqQkl0WlpDSjRfci0yaGtZdng0SVVMcVU1ajFwV2RjM2pwWE1qUS1lUl85YUxKSUtveEFtME1FYkZSWlBTWU1RTlpTY2pXN2dVZFBJRGxxX3VwenM4LUdRTWhVQW00Q2dVLVlXay1fZXZmSWhVRVh6ZUFtdFQzYWk1R3loYWRUYzNmUDBMOFlUT1ZsRXZWd251WFh6RV9aMjdkQnRpUXlrNk1LWmFjX09mX3ZpT2Rfazd5QktsbzZQVDZSMmJyMUtZTWVReVdpSnBrM3NLNlRacUdnWDEtWThDeFQ2cnowMjAzUmYyN3lwdnM4N1BKN0IyZHpqcHVIbGhFTjZTUDNpTGNYNThUYWphWm5GUnZPWGljTENtVUM5MGFtSFM5RjBxTXZBVFZpN25lUXRjZjlrWVVsdm15ODRaRFc5S0VKcUZrbV93NUtKRURNMGpUcW5DV0x5NVpkR29BaWsyd3J2Wm1vY1Z3OXE2a3YwWjA3dklEUVd2NzBTc0liR0pTZGpyYVh6d0hvZ2NVUTZHU3EyYXpub1cyS2k2RkZQaDlQZHJ5U3lKeEF6Si13TExvQTN1RSIsCiAgICAgICAgICJwcm90ZWN0ZWQiOiAiZXlKbWIzSnRZWFJNWlc1bmRHZ2lPakUxTVN3aVptOXliV0YwVkdGcGJDSTZJbVpSSWl3aWRHbHRaU0k2SWpJd01UVXRNVEl0TVRCVU1UZzZNekE2TVROYUluMCIKICAgICAgfQogICBdCn0="}}`)
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(data))
	require.Nil(t, err)
	resp, err := client.Do(req)
	require.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	log.Info("Succesfully added license")
}