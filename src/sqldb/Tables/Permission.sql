CREATE TABLE [dbo].[Permission]
(   
    [Type] VARCHAR(100) NOT NULL,
    [Email] VARCHAR(100) NOT NULL,
    [Created] DATETIME NOT NULL,
    [CreatedBy] VARCHAR(100) NOT NULL,
    CONSTRAINT PK_Permission PRIMARY KEY (Type, Email)
)