CREATE TABLE [dbo].[LegalConsultation]
(   
    [ItemId] UNIQUEIDENTIFIER NOT NULL,
    [Email] VARCHAR(100) NOT NULL,
    [Created] DATETIME NOT NULL,
    [CreatedBy] VARCHAR(100) NOT NULL
    CONSTRAINT PK_LegalConsultation PRIMARY KEY (ItemId, Email),
    CONSTRAINT FK_LegalConsultation_Items FOREIGN KEY (ItemId) REFERENCES Items(Id)
)