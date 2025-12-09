package com.example.curs3projectback.dto.dictionary;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class AddressResponse {
    Long id;
    String street;
    String house;
    String building;
    String apartment;
    String consumer;
}

