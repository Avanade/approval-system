﻿CREATE PROCEDURE PR_UserAccess_Update
(	   
        @Id int,
        @ProjectId int,
        @UserPrincipalName varchar(100),
        @IsActive BIT 

)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

UPDATE [dbo].[UserAccess]
   SET
        [ProjectId] = @ProjectId,
        [UserPrincipalName] =  @UserPrincipalName,
        [IsActive] = @IsActive,
        [Modified] = GETDATE(),
        [ModifiedBy] = @UserPrincipalName
    WHERE  
        [Id] = @Id
END