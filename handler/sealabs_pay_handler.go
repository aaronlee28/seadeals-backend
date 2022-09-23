package handler

//func (h *Handler) CreateSignature(ctx *gin.Context) {
//	value, _ := ctx.Get("payload")
//	json, _ := value.(*dto.TransactionDetailsReq)
//	transactionID := json.TransactionID
//
//	result, err := h.walletService.UserWalletData(userID)
//	if err != nil {
//		e := ctx.Error(err)
//		ctx.JSON(http.StatusBadRequest, e)
//		return
//	}
//	successResponse := dto.StatusOKResponse(result)
//	ctx.JSON(http.StatusOK, successResponse)
//}
