CREATE PROCEDURE [dbo].[PR_Items_Select]
(
	@Search VARCHAR(50) = '',
	@Offset INT = 0,
	@Filter INT = 10,
	@ItemType bit = NULL, -- NULL - ALL / 0 - REQUESTOR / 1 - APPROVER
	@User varchar(100) = NULL,
	@IsApproved int = 4,
	@RequestType varchar(100) = NULL,
	@Organization varchar(100) = NULL
)
AS
BEGIN
	SELECT
		dbo.UidToString(i.Id) AS ItemId
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
		, COUNT(*) AS Score
	  FROM [dbo].[Items] i
		INNER JOIN ApplicationModules am ON i.ApplicationModuleId = am.Id
		INNER JOIN Applications a ON am.ApplicationId = a.Id
		INNER JOIN ApprovalTypes t ON t.Id = am.ApprovalTypeId
		INNER JOIN ApprovalRequestApprovers ara ON i.Id = ara.ItemId
		INNER JOIN STRING_SPLIT(@Search, ' ') AS ss ON (i.Subject LIKE '%'+ss.value+'%' OR i.CreatedBy LIKE '%'+ss.value+'%')
	  WHERE
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
		) AND
		(
			@RequestType IS NULL OR
			(@RequestType IS NOT NULL AND i.ApplicationModuleId = @RequestType)
		) AND
		(
			@Organization IS NULL OR
			(@Organization IS NOT NULL AND i.Body LIKE '%'+@Organization+'%')
		)
	GROUP BY 
		i.Id,
		a.Id,
		a.Name,
		am.Id,
		am.Name,
		RespondedBy,
		Subject,
		Body,
		DateSent,
		DateResponded,
		IsApproved,
		ApproverRemarks,
		I.Created,
		T.ApproveText,
		T.RejectText,
		AllowReassign
	ORDER BY Score, I.Created DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END