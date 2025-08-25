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

func GetpaymentsprogramHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		marketplace_idVal, ok := args["marketplace_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: marketplace_id"), nil
		}
		marketplace_id, ok := marketplace_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: marketplace_id"), nil
		}
		payments_program_typeVal, ok := args["payments_program_type"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: payments_program_type"), nil
		}
		payments_program_type, ok := payments_program_typeVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: payments_program_type"), nil
		}
		url := fmt.Sprintf("%s/payments_program/%s/%s", cfg.BaseURL, marketplace_id, payments_program_type)
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
		var result models.PaymentsProgramResponse
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

func CreateGetpaymentsprogramTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_payments_program_marketplace_id_payments_program_type",
		mcp.WithDescription("<span class="tablenote"><b>Note:</b> This method is no longer applicable, as all seller accounts globally have been enabled for the new eBay payment and checkout flow.</span><br/><br/>This method returns whether or not the user is opted-in to the specified payments program. Sellers opt-in to payments programs by marketplace and you use the <b>marketplace_id</b> path parameter to specify the marketplace of the status flag you want returned."),
		mcp.WithString("marketplace_id", mcp.Required(), mcp.Description("This path parameter specifies the eBay marketplace of the payments program for which you want to retrieve the seller's status.")),
		mcp.WithString("payments_program_type", mcp.Required(), mcp.Description("This path parameter specifies the payments program whose status is returned by the call.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetpaymentsprogramHandler(cfg),
	}
}
