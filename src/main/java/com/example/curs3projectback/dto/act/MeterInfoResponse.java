package com.example.curs3projectback.dto.act;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class MeterInfoResponse {
    String brand;
    String serialNumber;
    Double readings;
}

