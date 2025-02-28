CREATE PROCEDURE [dbo].[PR_IPDisclosureRequest_Select_ByApprovalRequestId]
    @ApprovalRequestId [UNIQUEIDENTIFIER]
AS
BEGIN
    SELECT
        [Id],
        [RequestorName],
        [RequestorEmail],
        [IPTitle],
        [IPType],
        [IPDescription],
        [Reason],
        [IsApproved],
        [ApproverRemarks],
        [Created],
        [ResponseDate],
        [RespondedBy],
        [Created]
    FROM
        [dbo].[IPDisclosureRequest]
    WHERE
        [ApprovalRequestId] = @ApprovalRequestId
END

