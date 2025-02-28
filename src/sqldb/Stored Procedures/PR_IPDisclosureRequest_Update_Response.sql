CREATE PROCEDURE [dbo].[PR_IPDisclosureRequest_Update_Response]
  @ApprovalRequestId [UNIQUEIDENTIFIER],
  @IsApproved [BIT],
  @ApproverRemarks [VARCHAR](100),
  @ResponseDate [DATETIME],
  @RespondedBy [VARCHAR](100)
AS
BEGIN
  UPDATE [dbo].[IPDisclosureRequest]
  SET
    [IsApproved] = @IsApproved,
    [ApproverRemarks] = @ApproverRemarks,
    [ResponseDate] = @ResponseDate,
    [RespondedBy] = @RespondedBy
  WHERE [ApprovalRequestId] = @ApprovalRequestId
END
