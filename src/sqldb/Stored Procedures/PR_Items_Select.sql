CREATE PROCEDURE [dbo].[PR_Items_Select]
(
	@Search VARCHAR(50) = '',
	@Offset INT = 0,
	@Filter INT = 10,
	@ItemType bit = NULL, -- NULL - ALL / 0 - REQUESTOR / 1 - APPROVER
	@User varchar(100) = NULL,
	@IsApproved int = 4
)
AS
BEGIN
	SELECT
		DISTINCT dbo.UidToString(i.Id) AS ItemId
		, dbo.UidToString(a.Id) AS ApplicationId
		, a.Name AS Application
		, dbo.UidToString(am.Id) AS ApplicationModuleId
		, am.Name AS Module
		, i.RespondedBy
		, Subject
		, Body
		, DateSent
		, DateResponded
		, IsApproved
		, ApproverRemarks
		, I.Created
		, T.ApproveText
		, T.RejectText
	    , isnull(AllowReassign,'') as AllowReassign
	  FROM [dbo].[Items] i
		INNER JOIN ApplicationModules am ON i.ApplicationModuleId = am.Id
		INNER JOIN Applications a ON am.ApplicationId = a.Id
		INNER JOIN ApprovalTypes t ON t.Id = am.ApprovalTypeId
		INNER JOIN ApprovalRequestApprovers ara ON i.Id = ara.ItemId
	  WHERE
		Subject LIKE '%'+@Search+'%'
		AND
		(
			@ItemType IS NULL 
			OR 
			(@ItemType = 0 AND (@User IS NULL OR i.CreatedBy = @User)) -- Items by Requestor
			OR
			(@ItemType = 1 AND (@User IS NULL OR ara.ApproverEmail = @User)) -- Items by Approver
		) AND
		(
			(@IsApproved = 0 AND i.IsApproved IS NULL) OR -- Pending
			(@IsApproved = 1 AND i.IsApproved = 1) OR -- Approved
			(@IsApproved = 2 AND i.IsApproved = 0) OR -- Rejected
			(@IsApproved = 3 AND i.IsApproved IS NOT NULL) -- Closed (Rejected, Approved)
			-- If the value of IsApproved is 4 then select all
		)
	ORDER BY I.Created DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END