package com.example.curs3projectback.dto.photo;

import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class PhotoUploadRequest {
    private Long inspectionActId;
    private Long replacementActId;

    @Size(max = 1000, message = "note is too long")
    private String note;

    @AssertTrue(message = "Exactly one actId must be provided")
    public boolean isExactlyOneActSelected() {
        return (inspectionActId == null) ^ (replacementActId == null);
    }
}

