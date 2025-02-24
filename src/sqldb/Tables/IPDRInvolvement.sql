CREATE TABLE [dbo].[IPDRInvolvement]
(   
    [RequestId] INT NOT NULL,
    [InvolvementId] INT NOT NULL,
    CONSTRAINT FK_IPDRInvolvement_IPDisclosureRequest FOREIGN KEY (RequestId) REFERENCES IPDisclosureRequest(Id),
    CONSTRAINT FK_IPDRInvolvement_Involvement FOREIGN KEY (InvolvementId) REFERENCES Involvement(Id)
)