﻿CREATE PROCEDURE [dbo].[PR_Users_Select_ByUserPrincipalName]
(
	@UserPrincipalName varchar(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

SELECT 
		[UserPrincipalName],
		[GivenName],
		[SurName],
		[JobTitle],
		[GithubUser],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
  FROM 
		[dbo].[Users]
  WHERE  
	    [UserPrincipalName] = @UserPrincipalName 

END