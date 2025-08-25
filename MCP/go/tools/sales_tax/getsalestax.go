package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/account-api/mcp-server/config"
	"github.com/account-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetsalestaxHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		countryCodeVal, ok := args["countryCode"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: countryCode"), nil
		}
		countryCode, ok := countryCodeVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: countryCode"), nil
		}
		jurisdictionIdVal, ok := args["jurisdictionId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: jurisdictionId"), nil
		}
		jurisdictionId, ok := jurisdictionIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: jurisdictionId"), nil
		}
		url := fmt.Sprintf("%s/sales_tax/%s/%s", cfg.BaseURL, countryCode, jurisdictionId)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		if cfg.BearerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.BearerToken))
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.SalesTax
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetsalestaxTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_sales_tax_countryCode_jurisdictionId",
		mcp.WithDescription("This call gets the current sales tax table entry for a specific tax jurisdiction. Specify the jurisdiction to retrieve using the <b>countryCode</b> and <b>jurisdictionId</b> path parameters. All four response fields will be returned if a sales tax entry exists for the tax jurisdiction. Otherwise, the response will be returned as empty.<br/><br/><span class="tablenote"><b>Important!</b> In most US states and territories, eBay now 'collects and remits' sales tax, so sellers can no longer configure sales tax rates for these states/territories.</span>"),
		mcp.WithString("countryCode", mcp.Required(), mcp.Description("This path parameter specifies the two-letter <a href=\"https://www.iso.org/iso-3166-country-codes.html \" title=\"https://www.iso.org \" target=\"_blank\">ISO 3166</a> code for the country whose sales tax table you want to retrieve.")),
		mcp.WithString("jurisdictionId", mcp.Required(), mcp.Description("This path parameter specifies the ID of the sales tax jurisdiction for the tax table entry you want to retrieve. Retrieve valid jurisdiction IDs using <b>getSalesTaxJurisdictions</b> in the Metadata API.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetsalestaxHandler(cfg),
	}
}
