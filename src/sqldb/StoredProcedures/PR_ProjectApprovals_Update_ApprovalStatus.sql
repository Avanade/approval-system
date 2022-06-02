CREATE PROCEDURE [dbo].[PR_ProjectsApproval_Update_ApproverResponse]
(
  @ApprovalSystemGUID UNIQUEIDENTIFIER,
  @ApprovalStatusId INT,
  @ApprovalDescription varchar(500),
  @ApprovalDate DATETIME
)
AS
BEGIN
  -- SET NOCOUNT ON added to prevent extra result sets from
  -- interfering with SELECT statements.
  SET NOCOUNT ON

UPDATE
	[dbo].[ProjectApprovals]
  SET
    [ApprovalStatusId] = @ApprovalStatusId,
    [ApprovalDescription] = @ApprovalDescription,
    [ModifiedBy] = [ApproverUserPrincipalName],
    [Modified] = GETDATE(),
    [ApprovalDate] = convert(DATETIME, @ApprovalDate)
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID
END