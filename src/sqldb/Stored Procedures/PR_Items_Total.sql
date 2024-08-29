CREATE PROCEDURE [dbo].[PR_Items_Total]
(
	@Search VARCHAR(50) = '',
	@ItemType bit = NULL, -- NULL - ALL / 0 - REQUESTOR / 1 - APPROVER,
	@User varchar(100) = NULL,
	@IsApproved int = 4,
	@RequestType varchar(100) = NULL
)
AS
BEGIN
	SELECT
		COUNT(*) AS Total
	FROM (
		SELECT
			DISTINCT i.Id
		FROM [dbo].[Items] i
			INNER JOIN ApplicationModules am ON i.ApplicationModuleId = am.Id
			INNER JOIN Applications a ON am.ApplicationId = a.Id
			INNER JOIN ApprovalTypes t ON t.Id = am.ApprovalTypeId
			INNER JOIN ApprovalRequestApprovers ara ON i.Id = ara.ItemId
		WHERE
			Subject LIKE '%'+@Search+'%' AND
			(
				@ItemType IS NULL 
				OR 
				(@ItemType = 0 AND (@User IS NULL OR i.CreatedBy = @User))
				OR
				(
					@ItemType = 1 AND (
						@User IS NULL OR (
							ara.ApproverEmail = @User
						)
					)
				)
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
			)
		) AS Items
END