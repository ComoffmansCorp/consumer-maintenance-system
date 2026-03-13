package com.example.curs3projectback.service;

import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class CurrentUserService {

    private final UserRepository userRepository;

    public User getRequiredUser(Authentication authentication) {
        Long userId;
        try {
            userId = Long.parseLong(authentication.getName());
        } catch (NumberFormatException ex) {
            throw new ResourceNotFoundException("User not found");
        }
        return userRepository.findById(userId)
                .orElseThrow(() -> new ResourceNotFoundException("User not found"));
    }
}
