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

func GetfulfillmentpolicyHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		fulfillmentPolicyIdVal, ok := args["fulfillmentPolicyId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: fulfillmentPolicyId"), nil
		}
		fulfillmentPolicyId, ok := fulfillmentPolicyIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: fulfillmentPolicyId"), nil
		}
		url := fmt.Sprintf("%s/fulfillment_policy/%s", cfg.BaseURL, fulfillmentPolicyId)
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
		var result models.FulfillmentPolicy
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

func CreateGetfulfillmentpolicyTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_fulfillment_policy_fulfillmentPolicyId",
		mcp.WithDescription("This method retrieves the complete details of a fulfillment policy. Supply the ID of the policy you want to retrieve using the <b>fulfillmentPolicyId</b> path parameter."),
		mcp.WithString("fulfillmentPolicyId", mcp.Required(), mcp.Description("This path parameter specifies the ID of the fulfillment policy you want to retrieve.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetfulfillmentpolicyHandler(cfg),
	}
}
