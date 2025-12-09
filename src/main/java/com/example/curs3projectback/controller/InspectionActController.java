package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.act.InspectionActRequest;
import com.example.curs3projectback.dto.act.InspectionActResponse;
import com.example.curs3projectback.service.InspectionActService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/acts/inspection")
@RequiredArgsConstructor
public class InspectionActController {

    private final InspectionActService inspectionActService;

    @PostMapping
    public ResponseEntity<InspectionActResponse> create(@RequestBody InspectionActRequest request,
                                                        Authentication authentication) {
        return ResponseEntity.ok(inspectionActService.create(request, authentication));
    }

    @GetMapping
    public ResponseEntity<List<InspectionActResponse>> myActs(Authentication authentication) {
        return ResponseEntity.ok(inspectionActService.findMyActs(authentication));
    }
}

