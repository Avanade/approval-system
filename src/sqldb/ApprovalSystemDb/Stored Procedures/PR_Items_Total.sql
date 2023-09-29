CREATE PROCEDURE [dbo].[PR_Items_Total]
(
	@Search VARCHAR(50) = '',
	@ItemType bit = NULL, -- NULL - ALL / 0 - REQUESTOR / 1 - APPROVER,
	@User varchar(100) = NULL,
	@IsApproved int = -1 -- -1 - ALL / NULL - PENDING / 0 - REJECTED / 1 - APPROVED
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
							ara.ApproverEmail = @User OR
							i.ApproverEmail = @User -- OBSOLETE
						)
					)
				)
			) AND
			(
				(@IsApproved = -1 OR i.IsApproved = @IsApproved) 
				OR
				(@IsApproved IS NULL AND i.IsApproved IS NULL)
			)
		) AS Items
END