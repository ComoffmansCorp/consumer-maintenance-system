package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.act.ReplacementActRequest;
import com.example.curs3projectback.dto.act.ReplacementActResponse;
import com.example.curs3projectback.service.ReplacementActService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/acts/replacement")
@RequiredArgsConstructor
public class ReplacementActController {

    private final ReplacementActService replacementActService;

    @PostMapping
    public ResponseEntity<ReplacementActResponse> create(@RequestBody ReplacementActRequest request,
                                                         Authentication authentication) {
        return ResponseEntity.ok(replacementActService.create(request, authentication));
    }

    @GetMapping
    public ResponseEntity<List<ReplacementActResponse>> myActs(Authentication authentication) {
        return ResponseEntity.ok(replacementActService.findMyActs(authentication));
    }
}

