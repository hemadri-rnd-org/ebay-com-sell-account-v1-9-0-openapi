package main

import (
	"github.com/account-api/mcp-server/config"
	"github.com/account-api/mcp-server/models"
	tools_program "github.com/account-api/mcp-server/tools/program"
	tools_subscription "github.com/account-api/mcp-server/tools/subscription"
	tools_payment_policy "github.com/account-api/mcp-server/tools/payment_policy"
	tools_advertising_eligibility "github.com/account-api/mcp-server/tools/advertising_eligibility"
	tools_fulfillment_policy "github.com/account-api/mcp-server/tools/fulfillment_policy"
	tools_sales_tax "github.com/account-api/mcp-server/tools/sales_tax"
	tools_custom_policy "github.com/account-api/mcp-server/tools/custom_policy"
	tools_payments_program "github.com/account-api/mcp-server/tools/payments_program"
	tools_return_policy "github.com/account-api/mcp-server/tools/return_policy"
	tools_rate_table "github.com/account-api/mcp-server/tools/rate_table"
	tools_kyc "github.com/account-api/mcp-server/tools/kyc"
	tools_onboarding "github.com/account-api/mcp-server/tools/onboarding"
	tools_privilege "github.com/account-api/mcp-server/tools/privilege"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_program.CreateOptoutofprogramTool(cfg),
		tools_program.CreateGetoptedinprogramsTool(cfg),
		tools_subscription.CreateGetsubscriptionTool(cfg),
		tools_payment_policy.CreateGetpaymentpolicybynameTool(cfg),
		tools_advertising_eligibility.CreateGetadvertisingeligibilityTool(cfg),
		tools_payment_policy.CreateDeletepaymentpolicyTool(cfg),
		tools_payment_policy.CreateGetpaymentpolicyTool(cfg),
		tools_payment_policy.CreateUpdatepaymentpolicyTool(cfg),
		tools_fulfillment_policy.CreateGetfulfillmentpoliciesTool(cfg),
		tools_payment_policy.CreateGetpaymentpoliciesTool(cfg),
		tools_payment_policy.CreateCreatepaymentpolicyTool(cfg),
		tools_sales_tax.CreateDeletesalestaxTool(cfg),
		tools_sales_tax.CreateGetsalestaxTool(cfg),
		tools_sales_tax.CreateCreateorreplacesalestaxTool(cfg),
		tools_custom_policy.CreateGetcustompoliciesTool(cfg),
		tools_custom_policy.CreateCreatecustompolicyTool(cfg),
		tools_payments_program.CreateGetpaymentsprogramTool(cfg),
		tools_fulfillment_policy.CreateDeletefulfillmentpolicyTool(cfg),
		tools_fulfillment_policy.CreateGetfulfillmentpolicyTool(cfg),
		tools_fulfillment_policy.CreateUpdatefulfillmentpolicyTool(cfg),
		tools_return_policy.CreateGetreturnpoliciesTool(cfg),
		tools_return_policy.CreateCreatereturnpolicyTool(cfg),
		tools_return_policy.CreateDeletereturnpolicyTool(cfg),
		tools_return_policy.CreateGetreturnpolicyTool(cfg),
		tools_return_policy.CreateUpdatereturnpolicyTool(cfg),
		tools_fulfillment_policy.CreateGetfulfillmentpolicybynameTool(cfg),
		tools_rate_table.CreateGetratetablesTool(cfg),
		tools_kyc.CreateGetkycTool(cfg),
		tools_onboarding.CreateGetpaymentsprogramonboardingTool(cfg),
		tools_privilege.CreateGetprivilegesTool(cfg),
		tools_custom_policy.CreateGetcustompolicyTool(cfg),
		tools_custom_policy.CreateUpdatecustompolicyTool(cfg),
		tools_fulfillment_policy.CreateCreatefulfillmentpolicyTool(cfg),
		tools_program.CreateOptintoprogramTool(cfg),
		tools_return_policy.CreateGetreturnpolicybynameTool(cfg),
		tools_sales_tax.CreateGetsalestaxesTool(cfg),
	}
}
