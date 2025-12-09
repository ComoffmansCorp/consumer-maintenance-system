package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.auth.AuthResponse;
import com.example.curs3projectback.dto.auth.LoginRequest;
import com.example.curs3projectback.dto.auth.RegisterRequest;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final AuthenticationManager authenticationManager;
    private final JwtService jwtService;
    private final UserDetailsServiceImpl userDetailsService;

    public AuthResponse register(RegisterRequest request) {
        if (userRepository.existsByUsername(request.getUsername())) {
            throw new BadRequestException("Логин уже используется");
        }
        var user = User.builder()
                .username(request.getUsername())
                .password(passwordEncoder.encode(request.getPassword()))
                .fullName(request.getFullName())
                .role(request.getRole())
                .build();
        userRepository.save(user);
        var userDetails = userDetailsService.loadUserByUsername(user.getUsername());
        var token = jwtService.generateToken(userDetails);
        return new AuthResponse(token, user.getFullName(), user.getRole());
    }

    public AuthResponse login(LoginRequest request) {
        var authentication = new UsernamePasswordAuthenticationToken(request.getUsername(), request.getPassword());
        authenticationManager.authenticate(authentication);
        var user = userRepository.findByUsername(request.getUsername())
                .orElseThrow(() -> new BadRequestException("Неверный логин или пароль"));
        var userDetails = userDetailsService.loadUserByUsername(request.getUsername());
        var token = jwtService.generateToken(userDetails);
        return new AuthResponse(token, user.getFullName(), user.getRole());
    }
}

