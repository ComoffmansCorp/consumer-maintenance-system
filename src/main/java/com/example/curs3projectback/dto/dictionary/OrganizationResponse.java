package com.example.curs3projectback.dto.dictionary;

import com.example.curs3projectback.model.enums.OrganizationType;
import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class OrganizationResponse {
    Long id;
    String name;
    OrganizationType type;
    String description;
}

