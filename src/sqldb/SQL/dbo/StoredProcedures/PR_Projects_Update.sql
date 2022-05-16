﻿CREATE PROCEDURE PR_Projects_Update
(
	@ID Int,
	@Name varchar(50),
	@CoOwner varchar(100),
	@Description varchar(1000),
	@ConfirmAvaIP bit,
	@ConfirmEnabledSecurity bit,
	@Created varchar(100),
	@ModifiedBy varchar(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
UPDATE [dbo].[Projects]
   SET [Id] = @ID
      ,[Name] = @Name
      ,[CoOwner] = @CoOwner
      ,[Description] = @Description
      ,[ConfirmAvaIP] = @ConfirmAvaIP
      ,[ConfirmEnabledSecurity] =@ConfirmEnabledSecurity
      ,[Created] =	GETDATE()
      ,[CreatedBy] = @Created
      ,[Modified] =	GETDATE()
      ,[ModifiedBy] = @ModifiedBy
 WHERE  
	   [Id] =@Id or @Id IS NULL
END