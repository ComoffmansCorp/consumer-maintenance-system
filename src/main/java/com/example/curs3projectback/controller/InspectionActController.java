package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.act.InspectionActRequest;
import com.example.curs3projectback.dto.act.InspectionActResponse;
import com.example.curs3projectback.service.InspectionActService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/acts/inspection")
@RequiredArgsConstructor
public class InspectionActController {

    private final InspectionActService inspectionActService;

    @PostMapping
    public ResponseEntity<InspectionActResponse> create(@Valid @RequestBody InspectionActRequest request,
                                                        Authentication authentication) {
        log.info("HTTP POST /api/acts/inspection userId={} taskId={} addressId={}",
                authentication.getName(), request.getTaskId(), request.getAddressId());
        return ResponseEntity.ok(inspectionActService.create(request, authentication));
    }

    @GetMapping
    public ResponseEntity<List<InspectionActResponse>> myActs(Authentication authentication) {
        log.info("HTTP GET /api/acts/inspection userId={}", authentication.getName());
        return ResponseEntity.ok(inspectionActService.findMyActs(authentication));
    }
}

