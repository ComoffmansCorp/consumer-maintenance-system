package com.example.curs3projectback.dto.photo;

import lombok.Data;

@Data
public class PhotoUploadRequest {
    private Long inspectionActId;
    private Long replacementActId;
    private String note;
}

