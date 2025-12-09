package com.example.curs3projectback.dto.act;

import lombok.Data;

@Data
public class MeterInfoRequest {
    private String brand;
    private String serialNumber;
    private Double readings;
}

