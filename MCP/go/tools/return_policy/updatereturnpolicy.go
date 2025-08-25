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

func UpdatereturnpolicyHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		return_policy_idVal, ok := args["return_policy_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: return_policy_id"), nil
		}
		return_policy_id, ok := return_policy_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: return_policy_id"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody models.ReturnPolicyRequest
		
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
		url := fmt.Sprintf("%s/return_policy/%s", cfg.BaseURL, return_policy_id)
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
		var result models.SetReturnPolicyResponse
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

func CreateUpdatereturnpolicyTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_return_policy_return_policy_id",
		mcp.WithDescription("This method updates an existing return policy. Specify the policy you want to update using the <b>return_policy_id</b> path parameter. Supply a complete policy payload with the updates you want to make; this call overwrites the existing policy with the new details specified in the payload."),
		mcp.WithString("return_policy_id", mcp.Required(), mcp.Description("This path parameter specifies the ID of the return policy you want to update.")),
		mcp.WithString("returnShippingCostPayer", mcp.Description("Input parameter: This field indicates who is responsible for paying for the shipping charges for returned items. The field can be set to either <code>BUYER</code> or <code>SELLER</code>.  <br/><br/>Depending on the return policy and specifics of the return, either the buyer or the seller can be responsible for the return shipping costs. Note that the seller is always responsible for return shipping costs for SNAD-related issues.  <br/><br/>This field is conditionally required if <b>returnsAccepted</b> is set to <code>true</code>. For implementation help, refer to <a href='https://developer.ebay.com/api-docs/sell/account/types/api:ReturnShippingCostPayerEnum'>eBay API documentation</a>")),
		mcp.WithString("description", mcp.Description("Input parameter: A seller-defined description of the return business policy. This description is only for the seller's use, and is not exposed on any eBay pages.  <br/><br/><b>Max length</b>: 250")),
		mcp.WithBoolean("extendedHolidayReturnsOffered", mcp.Description("Input parameter: <p class=\"tablenote\"><span  style=\"color: #dd1e31;\"><b>Important!</b></span> This field is deprecated, since eBay no longer supports extended holiday returns. Any value supplied in this field is neither read nor returned.</p> ")),
		mcp.WithObject("internationalOverride", mcp.Description("Input parameter: This type defines the fields for a seller's international return policy. Sellers have the ability to set separate domestic and international return policies, but if an international return policy is not set, the same return policy settings specified for the domestic return policy are also used for returns for international buyers. ")),
		mcp.WithString("marketplaceId", mcp.Description("Input parameter: The ID of the eBay marketplace to which this return business policy applies.  For implementation help, refer to <a href='https://developer.ebay.com/api-docs/sell/account/types/ba:MarketplaceIdEnum'>eBay API documentation</a>")),
		mcp.WithBoolean("returnsAccepted", mcp.Description("Input parameter: If set to <code>true</code>, the seller accepts returns. <p><span class=\"tablenote\"><strong>Note:</strong>Top-Rated sellers must accept item returns and the <b>handlingTime</b> should be set to zero days or one day for a listing to receive a Top-Rated Plus badge on the View Item or search result pages. For more information on eBay's Top-Rated seller program, see <a href=\"http://pages.ebay.com/help/sell/top-rated.html \">Becoming a Top Rated Seller and qualifying for Top Rated Plus benefits</a>.</span></p>")),
		mcp.WithArray("categoryTypes", mcp.Description("Input parameter: This container indicates which category group that the return policy applies to.<br/><br/><span class=\"tablenote\"><b>Note</b>: Return business policies are not applicable to motor vehicle listings, so the <b>categoryTypes.name</b> value must be set to <code>ALL_EXCLUDING_MOTORS_VEHICLES</code> for return business policies.</span>")),
		mcp.WithString("refundMethod", mcp.Description("Input parameter: This value indicates the refund method that will be used by the seller for buyer returns.<p class=\"tablenote\"><span  style=\"color: #dd1e31;\"><b>Important!</b></span> If this field is not included in a return business policy, it will default to MONEY_BACK.</p> For implementation help, refer to <a href='https://developer.ebay.com/api-docs/sell/account/types/api:RefundMethodEnum'>eBay API documentation</a>")),
		mcp.WithString("returnInstructions", mcp.Description("Input parameter: This text-based field provides more details on seller-specified return instructions. <p class=\"tablenote\"><span  style=\"color: #dd1e31;\"><b>Important!</b></span> This field is no longer supported on many eBay marketplaces. To see if a marketplace and eBay category does support this field, call <a href=\"/api-docs/sell/metadata/resources/marketplace/methods/getReturnPolicies\">getReturnPolicies</a> method of the <b>Metadata API</b>. Then you will look for the <b>policyDescriptionEnabled</b> field with a value of <code>true</code> for the eBay category.</span></p><br/><b>Max length</b>: 5000 (8000 for DE)")),
		mcp.WithString("name", mcp.Description("Input parameter: A seller-defined name for this return business policy. Names must be unique for policies assigned to the same marketplace. <br/><br/><b>Max length</b>: 64")),
		mcp.WithString("restockingFeePercentage", mcp.Description("Input parameter: <p class=\"tablenote\"><span  style=\"color: #dd1e31;\"><b>Important!</b></span> This field is deprecated, since eBay no longer allows sellers to charge a restocking fee for buyer remorse returns. If this field is included, it is ignored.</p>")),
		mcp.WithString("returnMethod", mcp.Description("Input parameter: This field can be used if the seller is willing and able to offer a replacement item as an alternative to 'Money Back'. For implementation help, refer to <a href='https://developer.ebay.com/api-docs/sell/account/types/api:ReturnMethodEnum'>eBay API documentation</a>")),
		mcp.WithObject("returnPeriod", mcp.Description("Input parameter: A type used to specify a period of time using a specified time-measurement unit. Payment, return, and fulfillment business policies all use this type to specify time windows.<br/><br/>Whenever a container that uses this type is used in a request, both of these fields are required. Similarly, whenever a container that uses this type is returned in a response, both of these fields are always returned.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    UpdatereturnpolicyHandler(cfg),
	}
}
