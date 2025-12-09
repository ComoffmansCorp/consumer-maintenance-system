package com.example.curs3projectback.dto.photo;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class PhotoResponse {
    Long id;
    String url;
    String note;
}

