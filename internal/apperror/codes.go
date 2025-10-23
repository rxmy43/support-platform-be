package apperror

// Application-specific error codes for SaaS Restaurant Cashier App
// These codes are used as translation keys on the frontend

// ErrorCode represents application-specific error codes for i18n
type ErrorCode string

// Fallback Code -> use for unknown error code
const CodeUnknown ErrorCode = "unknown_error_code"

// Domain: API Response
const (
	CodeInternalServerError ErrorCode = "api.internal_server_error"
	CodeMethodNotAllowed    ErrorCode = "api.method_not_allowed"
	CodeServiceUnavailable  ErrorCode = "api.code_service_unavailable"
)

// Domain: Field Validations
const (
	// General field validation
	CodeFieldRequired            ErrorCode = "validation.field_required"
	CodeFieldTooShort            ErrorCode = "validation.field_too_short"
	CodeFieldTooLong             ErrorCode = "validation.field_too_long"
	CodeFieldInvalidFormat       ErrorCode = "validation.field_invalid_format"
	CodeFieldInvalidType         ErrorCode = "validation.field_invalid_type"
	CodeFieldOutOfRange          ErrorCode = "validation.field_out_of_range"
	CodeFieldDuplicate           ErrorCode = "validation.field_duplicate"
	CodeFieldReadOnly            ErrorCode = "validation.field_read_only"
	CodeFieldImmutable           ErrorCode = "validation.field_immutable"
	CodeInvalidRequestJSONFormat ErrorCode = "validation.invalid_request_json_format"
	CodeResourceNotFound         ErrorCode = "validation.resource_not_found"
	CodeInvalidCredentials       ErrorCode = "validation.invalid_credentials"

	// Validator error (500 not 422)
	CodeInvalidValidatorTarget ErrorCode = "validation.invalid_validator_target"
	CodeInvalidRequestStruct   ErrorCode = "validation.invalid_request_struct"

	// Specific field types
	CodeEmailInvalid       ErrorCode = "validation.email_invalid"
	CodeEmailDomainInvalid ErrorCode = "validation.email_domain_invalid"
	CodeEmailTaken         ErrorCode = "validation.email_taken"
	CodePhoneTaken         ErrorCode = "validation.phone_taken"
	CodePhoneInvalid       ErrorCode = "validation.phone_invalid"
	CodePhoneFormatInvalid ErrorCode = "validation.phone_format_invalid"
	CodeURLInvalid         ErrorCode = "validation.url_invalid"
	CodeDateInvalid        ErrorCode = "validation.date_invalid"
	CodeDateInPast         ErrorCode = "validation.date_in_past"
	CodeDateInFuture       ErrorCode = "validation.date_in_future"
	CodeDateTimeInvalid    ErrorCode = "validation.datetime_invalid"
	CodeTimeInvalid        ErrorCode = "validation.time_invalid"

	// Numeric validations
	CodeNumberTooSmall        ErrorCode = "validation.number_too_small"
	CodeNumberTooLarge        ErrorCode = "validation.number_too_large"
	CodeNumberInvalid         ErrorCode = "validation.number_invalid"
	CodeDecimalPlacesExceeded ErrorCode = "validation.decimal_places_exceeded"
	CodeNegativeNotAllowed    ErrorCode = "validation.negative_not_allowed"
	CodeZeroNotAllowed        ErrorCode = "validation.zero_not_allowed"
	CodePositiveRequired      ErrorCode = "validation.positive_required"

	// String validations
	CodeStringFormatInvalid  ErrorCode = "validation.string_format_invalid"
	CodeInvalidCharacters    ErrorCode = "validation.invalid_characters"
	CodeSpecialCharsRequired ErrorCode = "validation.special_chars_required"
	CodeUppercaseRequired    ErrorCode = "validation.uppercase_required"
	CodeLowercaseRequired    ErrorCode = "validation.lowercase_required"
	CodeDigitsRequired       ErrorCode = "validation.digits_required"
	CodeNoWhitespaceAllowed  ErrorCode = "validation.no_whitespace_allowed"
	CodeWhitespaceRequired   ErrorCode = "validation.whitespace_required"

	// Password validations
	CodePasswordTooWeak      ErrorCode = "validation.password_too_weak"
	CodePasswordNoMatch      ErrorCode = "validation.password_no_match"
	CodePasswordMin6Chars    ErrorCode = "validation.password_min_6_chars"
	CodePasswordCommon       ErrorCode = "validation.password_common"
	CodePasswordSequential   ErrorCode = "validation.password_sequential"
	CodePasswordPersonalInfo ErrorCode = "validation.password_personal_info"

	// Business-specific field validations
	CodePINInvalid           ErrorCode = "validation.pin_invalid"
	CodePINTooShort          ErrorCode = "validation.pin_too_short"
	CodeStaffCodeInvalid     ErrorCode = "validation.staff_code_invalid"
	CodeTableNumberInvalid   ErrorCode = "validation.table_number_invalid"
	CodeOrderNumberInvalid   ErrorCode = "validation.order_number_invalid"
	CodeReceiptNumberInvalid ErrorCode = "validation.receipt_number_invalid"

	// Price and money validations
	CodePriceInvalid         ErrorCode = "validation.price_invalid"
	CodePriceTooLow          ErrorCode = "validation.price_too_low"
	CodePriceTooHigh         ErrorCode = "validation.price_too_high"
	CodeCurrencyInvalid      ErrorCode = "validation.currency_invalid"
	CodeTaxRateInvalid       ErrorCode = "validation.tax_rate_invalid"
	CodeDiscountInvalid      ErrorCode = "validation.discount_invalid"
	CodeTipPercentageInvalid ErrorCode = "validation.tip_percentage_invalid"

	// Quantity validations
	CodeQuantityInvalid       ErrorCode = "validation.quantity_invalid"
	CodeQuantityZero          ErrorCode = "validation.quantity_zero"
	CodeQuantityNegative      ErrorCode = "validation.quantity_negative"
	CodeQuantityDecimal       ErrorCode = "validation.quantity_decimal"
	CodeQuantityWholeRequired ErrorCode = "validation.quantity_whole_required"

	// Time validations
	CodeTimeRangeInvalid      ErrorCode = "validation.time_range_invalid"
	CodeOperatingHoursInvalid ErrorCode = "validation.operating_hours_invalid"
	CodeDurationInvalid       ErrorCode = "validation.duration_invalid"
	CodeDateRangeInvalid      ErrorCode = "validation.date_range_invalid"

	// Selection validations
	CodeSelectionRequired        ErrorCode = "validation.selection_required"
	CodeSelectionInvalid         ErrorCode = "validation.selection_invalid"
	CodeMultipleSelectionInvalid ErrorCode = "validation.multiple_selection_invalid"
	CodeOptionNotAvailable       ErrorCode = "validation.option_not_available"

	// File validations
	CodeFileTooLarge            ErrorCode = "validation.file_too_large"
	CodeFileTypeInvalid         ErrorCode = "validation.file_type_invalid"
	CodeFileCorrupted           ErrorCode = "validation.file_corrupted"
	CodeFileUploadFailed        ErrorCode = "validation.file_upload_failed"
	CodeImageDimensionsInvalid  ErrorCode = "validation.image_dimensions_invalid"
	CodeImageAspectRatioInvalid ErrorCode = "validation.image_aspect_ratio_invalid"

	// Location validations
	CodeAddressInvalid     ErrorCode = "validation.address_invalid"
	CodePostalCodeInvalid  ErrorCode = "validation.postal_code_invalid"
	CodeCityInvalid        ErrorCode = "validation.city_invalid"
	CodeStateInvalid       ErrorCode = "validation.state_invalid"
	CodeCountryInvalid     ErrorCode = "validation.country_invalid"
	CodeCoordinatesInvalid ErrorCode = "validation.coordinates_invalid"

	// Payment validations
	CodeCardNumberInvalid     ErrorCode = "validation.card_number_invalid"
	CodeCVVInvalid            ErrorCode = "validation.cvv_invalid"
	CodeCardholderNameInvalid ErrorCode = "validation.cardholder_name_invalid"
	CodePaymentMethodInvalid  ErrorCode = "validation.payment_method_invalid"

	// Restaurant-specific validations
	CodeTableCapacityInvalid ErrorCode = "validation.table_capacity_invalid"
	CodeSeatingNumberInvalid ErrorCode = "validation.seating_number_invalid"
	CodeMenuCategoryInvalid  ErrorCode = "validation.menu_category_invalid"
	CodeMenuItemCodeInvalid  ErrorCode = "validation.menu_item_code_invalid"
	CodeModifierGroupInvalid ErrorCode = "validation.modifier_group_invalid"
	CodeCookingTimeInvalid   ErrorCode = "validation.cooking_time_invalid"
	CodePrepTimeInvalid      ErrorCode = "validation.prep_time_invalid"

	// Inventory validations
	CodeSKUInvalid           ErrorCode = "validation.sku_invalid"
	CodeBarcodeInvalid       ErrorCode = "validation.barcode_invalid"
	CodeUnitOfMeasureInvalid ErrorCode = "validation.unit_of_measure_invalid"
	CodeSupplierCodeInvalid  ErrorCode = "validation.supplier_code_invalid"
	CodeBatchNumberInvalid   ErrorCode = "validation.batch_number_invalid"
	CodeExpiryDateInvalid    ErrorCode = "validation.expiry_date_invalid"

	// Customer validations
	CodeCustomerNameInvalid      ErrorCode = "validation.customer_name_invalid"
	CodeLoyaltyNumberInvalid     ErrorCode = "validation.loyalty_number_invalid"
	CodeAllergyInfoInvalid       ErrorCode = "validation.allergy_info_invalid"
	CodeDietaryPreferenceInvalid ErrorCode = "validation.dietary_preference_invalid"
)

// Domain: Authentication & Session Management
const (
	// Session, Token, and Header
	CodeTokenExpired     ErrorCode = "auth.token_expired"
	CodeTokenInvalid     ErrorCode = "auth.token_invalid"
	CodeTokenNotFound    ErrorCode = "auth.token_not_found"
	CodeMissingXTenantID ErrorCode = "auth.missing_x_tenant_id"

	// Staff authentication
	CodeStaffNotFound       ErrorCode = "auth.staff_not_found"
	CodeInvalidPin          ErrorCode = "auth.invalid_pin"
	CodeStaffInactive       ErrorCode = "auth.staff_inactive"
	CodeStaffSuspended      ErrorCode = "auth.staff_suspended"
	CodeShiftNotStarted     ErrorCode = "auth.shift_not_started"
	CodeShiftAlreadyStarted ErrorCode = "auth.shift_already_started"
	CodeShiftEnded          ErrorCode = "auth.shift_ended"

	// Session and permissions
	CodeUnauthorizedOperation ErrorCode = "auth.unauthorized_operation"
	CodeManagerRequired       ErrorCode = "auth.manager_required"
	CodeKitchenAccessRequired ErrorCode = "auth.kitchen_access_required"
	CodeCashierAccessRequired ErrorCode = "auth.cashier_access_required"
	CodeAdminAccessRequired   ErrorCode = "auth.admin_access_required"
)

// Domain: Restaurant & Location Management
const (
	CodeRestaurantNotFound  ErrorCode = "restaurant.not_found"
	CodeRestaurantDisabled  ErrorCode = "restaurant.disabled"
	CodeRestaurantSuspended ErrorCode = "restaurant.suspended"
	CodeLocationNotFound    ErrorCode = "restaurant.location_not_found"
	CodeLocationClosed      ErrorCode = "restaurant.location_closed"
	CodeLocationOffline     ErrorCode = "restaurant.location_offline"

	// Operating hours
	CodeOutsideOperatingHours ErrorCode = "restaurant.outside_operating_hours"
	CodeSpecialHoursActive    ErrorCode = "restaurant.special_hours_active"
	CodeHolidaySchedule       ErrorCode = "restaurant.holiday_schedule"
)

// Domain: Table & Floor Management
const (
	CodeTableNotFound         ErrorCode = "table.not_found"
	CodeTableOccupied         ErrorCode = "table.occupied"
	CodeTableReserved         ErrorCode = "table.reserved"
	CodeTableNotAvailable     ErrorCode = "table.not_available"
	CodeTableCapacityExceeded ErrorCode = "table.capacity_exceeded"
	CodeTableMergeNotAllowed  ErrorCode = "table.merge_not_allowed"
	CodeTableSplitNotAllowed  ErrorCode = "table.split_not_allowed"

	// Floor plan
	CodeFloorPlanNotFound  ErrorCode = "table.floor_plan_not_found"
	CodeSectionNotFound    ErrorCode = "table.section_not_found"
	CodeInvalidTableLayout ErrorCode = "table.invalid_layout"
)

// Domain: Order Management
const (
	CodeOrderNotFound      ErrorCode = "order.not_found"
	CodeOrderAlreadyExists ErrorCode = "order.already_exists"
	CodeOrderClosed        ErrorCode = "order.closed"
	CodeOrderVoided        ErrorCode = "order.voided"
	CodeOrderInProgress    ErrorCode = "order.in_progress"
	CodeOrderAlreadyPaid   ErrorCode = "order.already_paid"
	CodeOrderNotPaid       ErrorCode = "order.not_paid"

	// Order operations
	CodeModifyClosedOrder   ErrorCode = "order.modify_closed"
	CodeVoidAfterPayment    ErrorCode = "order.void_after_payment"
	CodeSplitOrderFailed    ErrorCode = "order.split_failed"
	CodeMergeOrderFailed    ErrorCode = "order.merge_failed"
	CodeTransferOrderFailed ErrorCode = "order.transfer_failed"
	CodeDiscountNotAllowed  ErrorCode = "order.discount_not_allowed"

	// Order timing
	CodeOrderTooOld          ErrorCode = "order.too_old"
	CodeModifyAfterTimeLimit ErrorCode = "order.modify_time_limit"
)

// Domain: Menu & Items
const (
	CodeMenuItemNotFound     ErrorCode = "menu.item_not_found"
	CodeMenuItemUnavailable  ErrorCode = "menu.item_unavailable"
	CodeMenuCategoryNotFound ErrorCode = "menu.category_not_found"
	CodeMenuSectionNotFound  ErrorCode = "menu.section_not_found"
	CodeMenuModifierNotFound ErrorCode = "menu.modifier_not_found"

	// Menu operations
	CodeItemOutOfStock        ErrorCode = "menu.out_of_stock"
	CodeItemSeasonal          ErrorCode = "menu.seasonal"
	CodeItemDaypartRestricted ErrorCode = "menu.daypart_restricted"
	CodeItemPriceChanged      ErrorCode = "menu.price_changed"
	CodeModifierRequired      ErrorCode = "menu.modifier_required"
	CodeModifierConflict      ErrorCode = "menu.modifier_conflict"
	CodeComboInvalid          ErrorCode = "menu.combo_invalid"

	// Menu management
	CodeMenuVersionMismatch  ErrorCode = "menu.version_mismatch"
	CodeMenuNotPublished     ErrorCode = "menu.not_published"
	CodeMenuUpdateInProgress ErrorCode = "menu.update_in_progress"
)

// Domain: Order Items & Modifications
const (
	CodeInvalidItemQuantity       ErrorCode = "order_item.invalid_quantity"
	CodeItemAlreadyAdded          ErrorCode = "order_item.already_added"
	CodeItemModificationFailed    ErrorCode = "order_item.modification_failed"
	CodeSpecialInstructionInvalid ErrorCode = "order_item.special_instruction_invalid"
	CodeCookingInstructionInvalid ErrorCode = "order_item.cooking_instruction_invalid"

	// Kitchen instructions
	CodeAllergyWarningRequired  ErrorCode = "order_item.allergy_warning_required"
	CodePreparationTimeExceeded ErrorCode = "order_item.preparation_time_exceeded"
	CodeItemHoldFailed          ErrorCode = "order_item.hold_failed"
	CodeItemFireFailed          ErrorCode = "order_item.fire_failed"
)

// Domain: Payment Processing
const (
	CodePaymentNotFound         ErrorCode = "payment.not_found"
	CodePaymentFailed           ErrorCode = "payment.failed"
	CodePaymentDeclined         ErrorCode = "payment.declined"
	CodePaymentProcessorDown    ErrorCode = "payment.processor_down"
	CodePaymentAlreadyRefunded  ErrorCode = "payment.already_refunded"
	CodePartialRefundNotAllowed ErrorCode = "payment.partial_refund_not_allowed"

	// Payment methods
	CodeCardDeclined           ErrorCode = "payment.card_declined"
	CodeCardExpired            ErrorCode = "payment.card_expired"
	CodeInvalidCard            ErrorCode = "payment.invalid_card"
	CodeCashPaymentExactAmount ErrorCode = "payment.cash_exact_amount"
	CodeMobilePaymentFailed    ErrorCode = "payment.mobile_payment_failed"
	CodeGiftCardInsufficient   ErrorCode = "payment.gift_card_insufficient"
	CodeGiftCardExpired        ErrorCode = "payment.gift_card_expired"

	// Tip processing
	CodeTipTooHigh          ErrorCode = "payment.tip_too_high"
	CodeTipBelowMinimum     ErrorCode = "payment.tip_below_minimum"
	CodeTipAlreadyProcessed ErrorCode = "payment.tip_already_processed"
)

// Domain: Discounts & Promotions
const (
	CodeDiscountNotFound         ErrorCode = "discount.not_found"
	CodeDiscountExpired          ErrorCode = "discount.expired"
	CodeDiscountNotStarted       ErrorCode = "discount.not_started"
	CodeDiscountUsageLimit       ErrorCode = "discount.usage_limit"
	CodeDiscountMinimumAmount    ErrorCode = "discount.minimum_amount"
	CodeDiscountItemExcluded     ErrorCode = "discount.item_excluded"
	CodeDiscountCategoryExcluded ErrorCode = "discount.category_excluded"
	CodeDiscountAlreadyApplied   ErrorCode = "discount.already_applied"
	CodeDiscountComboInvalid     ErrorCode = "discount.combo_invalid"

	// Coupons and vouchers
	CodeCouponInvalid             ErrorCode = "discount.coupon_invalid"
	CodeCouponAlreadyRedeemed     ErrorCode = "discount.coupon_redeemed"
	CodeLoyaltyPointsInsufficient ErrorCode = "discount.loyalty_points_insufficient"
)

// Domain: Customer Management
const (
	CodeCustomerNotFound        ErrorCode = "customer.not_found"
	CodeCustomerExists          ErrorCode = "customer.already_exists"
	CodeCustomerVIPOnly         ErrorCode = "customer.vip_only"
	CodeCustomerBlacklisted     ErrorCode = "customer.blacklisted"
	CodeLoyaltyTierInsufficient ErrorCode = "customer.loyalty_tier_insufficient"

	// Customer preferences
	CodeAllergyConflict    ErrorCode = "customer.allergy_conflict"
	CodeDietaryRestriction ErrorCode = "customer.dietary_restriction"
	CodePreferenceConflict ErrorCode = "customer.preference_conflict"
)

// Domain: Kitchen & Production
const (
	CodeKitchenStationNotFound ErrorCode = "kitchen.station_not_found"
	CodeKitchenStationOffline  ErrorCode = "kitchen.station_offline"
	CodeKitchenOverloaded      ErrorCode = "kitchen.overloaded"
	CodePrepStationClosed      ErrorCode = "kitchen.prep_station_closed"

	// Food preparation
	CodeItemPrepTimeExceeded ErrorCode = "kitchen.prep_time_exceeded"
	CodeIngredientShortage   ErrorCode = "kitchen.ingredient_shortage"
	CodeEquipmentMalfunction ErrorCode = "kitchen.equipment_malfunction"
	CodeQualityCheckFailed   ErrorCode = "kitchen.quality_check_failed"

	// Order timing
	CodeOrderRush             ErrorCode = "kitchen.order_rush"
	CodeEstimatedTimeExceeded ErrorCode = "kitchen.estimated_time_exceeded"
)

// Domain: Inventory & Stock Management
const (
	CodeInventoryItemNotFound   ErrorCode = "inventory.item_not_found"
	CodeInsufficientStock       ErrorCode = "inventory.insufficient_stock"
	CodeStockCountMismatch      ErrorCode = "inventory.count_mismatch"
	CodeInventoryUpdateConflict ErrorCode = "inventory.update_conflict"
	CodeSupplierOutOfStock      ErrorCode = "inventory.supplier_out_of_stock"
	CodeWasteTrackingRequired   ErrorCode = "inventory.waste_tracking_required"

	// Inventory operations
	CodeBatchExpired        ErrorCode = "inventory.batch_expired"
	CodeBatchRecall         ErrorCode = "inventory.batch_recall"
	CodeParLevelExceeded    ErrorCode = "inventory.par_level_exceeded"
	CodeReorderPointReached ErrorCode = "inventory.reorder_point_reached"
)

// Domain: Reporting & Analytics
const (
	CodeReportGenerationFailed ErrorCode = "report.generation_failed"
	CodeReportDataIncomplete   ErrorCode = "report.data_incomplete"
	CodeReportPeriodInvalid    ErrorCode = "report.period_invalid"
	CodeSalesDataCorrupted     ErrorCode = "report.sales_data_corrupted"
	CodeLaborReportUnavailable ErrorCode = "report.labor_unavailable"
	CodeInventoryReportPending ErrorCode = "report.inventory_pending"
)

// Domain: System & Integration
const (
	// Hardware integration
	CodePrinterOffline         ErrorCode = "system.printer_offline"
	CodePrinterPaperOut        ErrorCode = "system.printer_paper_out"
	CodeScannerOffline         ErrorCode = "system.scanner_offline"
	CodeScaleOffline           ErrorCode = "system.scale_offline"
	CodePaymentTerminalOffline ErrorCode = "system.payment_terminal_offline"
	CodeKDSOffline             ErrorCode = "system.kds_offline"

	// External services
	CodeAccountingSyncFailed  ErrorCode = "system.accounting_sync_failed"
	CodeOnlineOrderingDown    ErrorCode = "system.online_ordering_down"
	CodeDeliveryServiceDown   ErrorCode = "system.delivery_service_down"
	CodeReservationSystemDown ErrorCode = "system.reservation_system_down"

	// Backoffice systems
	CodeBackofficeUnavailable ErrorCode = "system.backoffice_unavailable"
	CodeDataSyncFailed        ErrorCode = "system.data_sync_failed"
	CodeBackupInProgress      ErrorCode = "system.backup_in_progress"
)

// Domain: Tax & Compliance
const (
	CodeTaxCalculationFailed   ErrorCode = "tax.calculation_failed"
	CodeTaxRateNotFound        ErrorCode = "tax.rate_not_found"
	CodeTaxExemptInvalid       ErrorCode = "tax.exempt_invalid"
	CodeTaxCertificateExpired  ErrorCode = "tax.certificate_expired"
	CodeComplianceCheckFailed  ErrorCode = "tax.compliance_check_failed"
	CodeSalesTaxReportingError ErrorCode = "tax.reporting_error"
)

// Domain: Delivery & Takeout
const (
	CodeDeliveryZoneUnavailable ErrorCode = "delivery.zone_unavailable"
	CodeDeliveryTimeExceeded    ErrorCode = "delivery.time_exceeded"
	CodeDriverUnavailable       ErrorCode = "delivery.driver_unavailable"
	CodeTakeoutReadyTimeInvalid ErrorCode = "delivery.ready_time_invalid"
	CodePackagingUnavailable    ErrorCode = "delivery.packaging_unavailable"
)
