package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/account-api/mcp-server/config"
	"github.com/account-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func UpdatepaymentpolicyHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		payment_policy_idVal, ok := args["payment_policy_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: payment_policy_id"), nil
		}
		payment_policy_id, ok := payment_policy_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: payment_policy_id"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody models.PaymentPolicyRequest
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/payment_policy/%s", cfg.BaseURL, payment_policy_id)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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
		var result models.SetPaymentPolicyResponse
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

func CreateUpdatepaymentpolicyTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_payment_policy_payment_policy_id",
		mcp.WithDescription("This method updates an existing payment policy. Specify the policy you want to update using the <b>payment_policy_id</b> path parameter. Supply a complete policy payload with the updates you want to make; this call overwrites the existing policy with the new details specified in the payload."),
		mcp.WithString("payment_policy_id", mcp.Required(), mcp.Description("This path parameter specifies the ID of the payment policy you want to update.")),
		mcp.WithArray("categoryTypes", mcp.Description("Input parameter: This container is used to specify whether the payment business policy applies to motor vehicle listings, or if it applies to non-motor vehicle listings.")),
		mcp.WithString("description", mcp.Description("Input parameter: A seller-defined description of the payment business policy. This description is only for the seller's use, and is not exposed on any eBay pages.  <br/><br/><b>Max length</b>: 250")),
		mcp.WithString("name", mcp.Description("Input parameter: A seller-defined name for this payment business policy. Names must be unique for policies assigned to the same marketplace.<br /><br /><b>Max length:</b> 64")),
		mcp.WithObject("deposit", mcp.Description("Input parameter: This type is used to specify/indicate that an initial deposit is required for a motor vehicle listing.")),
		mcp.WithObject("fullPaymentDueIn", mcp.Description("Input parameter: A type used to specify a period of time using a specified time-measurement unit. Payment, return, and fulfillment business policies all use this type to specify time windows.<br/><br/>Whenever a container that uses this type is used in a request, both of these fields are required. Similarly, whenever a container that uses this type is returned in a response, both of these fields are always returned.")),
		mcp.WithBoolean("immediatePay", mcp.Description("Input parameter: This field should be included and set to <code>true</code> if the seller wants to require immediate payment from the buyer for: <ul><li>A fixed-price item</li><li>An auction item where the buyer is using the 'Buy it Now' option</li><li>A deposit for a motor vehicle listing</li></ul><br /><b>Default:</b> False")),
		mcp.WithString("marketplaceId", mcp.Description("Input parameter: The ID of the eBay marketplace to which this payment business policy applies. For implementation help, refer to <a href='https://developer.ebay.com/api-docs/sell/account/types/ba:MarketplaceIdEnum'>eBay API documentation</a>")),
		mcp.WithString("paymentInstructions", mcp.Description("Input parameter: <p class=\"tablenote\"><b>Note:</b> DO NOT USE THIS FIELD. Payment instructions are no longer supported by payment business policies.</p>A free-form string field that allows sellers to add detailed payment instructions to their listings.")),
		mcp.WithArray("paymentMethods", mcp.Description("Input parameter: <p class=\"tablenote\"><b>Note:</b> This field applies only when the seller needs to specify one or more offline payment methods. eBay now manages the electronic payment options available to buyers to pay for the item.</p>This array is used to specify one or more offline payment methods that will be accepted for payment that occurs off of eBay's platform.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    UpdatepaymentpolicyHandler(cfg),
	}
}
