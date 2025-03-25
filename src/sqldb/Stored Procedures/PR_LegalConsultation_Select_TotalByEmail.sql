CREATE PROCEDURE  [dbo].[PR_LegalConsultation_Select_TotalByEmail]
(
    @Email VARCHAR(100),
    @IsApproved INT = 4
)
AS
BEGIN   
    SELECT 
        COUNT(*) AS Total
    FROM [dbo].[LegalConsultation] AS [LC]
    INNER JOIN [dbo].[Items] AS [I] ON [I].[Id] = [LC].[ItemId]
    WHERE [LC].[Email] = @Email AND
		(
			(@IsApproved = 0 AND i.IsApproved IS NULL) OR -- Pending
			(@IsApproved = 1 AND i.IsApproved = 1) OR -- Approved
			(@IsApproved = 2 AND i.IsApproved = 0) OR -- Rejected
			(@IsApproved = 3 AND i.IsApproved IS NOT NULL) -- Closed (Rejected, Approved)
			-- If the value of IsApproved is 4 then select all
		)
END