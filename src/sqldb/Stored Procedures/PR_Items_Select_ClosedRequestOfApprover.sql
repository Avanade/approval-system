CREATE PROCEDURE [dbo].[PR_Items_Select_ClosedRequestOfApprover]
(
    @Search VARCHAR(50) = '',
    @Approver VARCHAR(100) = NULL,
    @Offset INT = 0,
    @Filter INT = 10,
    @RequestType varchar(100) = NULL,
    @Organization varchar(100) = NULL,
    @IsApproved BIT = NULL -- NULL - Closed (Approved, Rejected) ; 1 - Approved ; 0 - Rejected
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
		, i.CreatedBy AS RequestedBy
	    , isnull(AllowReassign,'') as AllowReassign
		, COUNT(*) AS Score
	  FROM [dbo].[Items] i
		INNER JOIN ApplicationModules am ON i.ApplicationModuleId = am.Id
		INNER JOIN Applications a ON am.ApplicationId = a.Id
		INNER JOIN ApprovalTypes t ON t.Id = am.ApprovalTypeId
		INNER JOIN ApprovalRequestApprovers ara ON i.Id = ara.ItemId
		INNER JOIN STRING_SPLIT(@Search, ' ') AS ss ON (i.Subject LIKE '%'+ss.value+'%' OR i.CreatedBy LIKE '%'+ss.value+'%')
    WHERE
        ara.ApproverEmail = @Approver AND
        (
            (@IsApproved IS NULL AND i.IsApproved IS NOT NULL) OR -- Closed (Rejected, Approved)
            (@IsApproved IS NOT NULL AND i.IsApproved = @IsApproved)
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
		AllowReassign,
		i.CreatedBy
	ORDER BY Score, I.Created DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END