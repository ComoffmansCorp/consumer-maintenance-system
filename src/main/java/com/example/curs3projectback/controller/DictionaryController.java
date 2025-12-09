package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.dictionary.AddressResponse;
import com.example.curs3projectback.dto.dictionary.OrganizationResponse;
import com.example.curs3projectback.service.DictionaryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/dictionaries")
@RequiredArgsConstructor
public class DictionaryController {

    private final DictionaryService dictionaryService;

    @GetMapping("/addresses")
    public ResponseEntity<List<AddressResponse>> addresses(@RequestParam(defaultValue = "") String q) {
        return ResponseEntity.ok(dictionaryService.findAddresses(q));
    }

    @GetMapping("/organizations")
    public ResponseEntity<List<OrganizationResponse>> organizations(@RequestParam(defaultValue = "") String q) {
        return ResponseEntity.ok(dictionaryService.findOrganizations(q));
    }
}

