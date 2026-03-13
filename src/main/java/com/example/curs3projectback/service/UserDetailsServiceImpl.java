package com.example.curs3projectback.service;

import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class UserDetailsServiceImpl implements UserDetailsService {

    private static final String LOGIN_SEPARATOR = "::";

    private final UserRepository userRepository;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        var user = isNumeric(username)
                ? userRepository.findById(Long.parseLong(username))
                : username.contains(LOGIN_SEPARATOR)
                ? userRepository.findByTenant_CodeIgnoreCaseAndUsernameIgnoreCase(
                        username.substring(0, username.indexOf(LOGIN_SEPARATOR)),
                        username.substring(username.indexOf(LOGIN_SEPARATOR) + LOGIN_SEPARATOR.length())
                )
                : userRepository.findByTenantIsNullAndUsernameIgnoreCase(username);

        var resolvedUser = user.orElseThrow(() -> new UsernameNotFoundException("User not found"));
        return new User(
                String.valueOf(resolvedUser.getId()),
                resolvedUser.getPassword(),
                List.of(new SimpleGrantedAuthority("ROLE_" + resolvedUser.getRole().name()))
        );
    }

    private boolean isNumeric(String value) {
        try {
            Long.parseLong(value);
            return true;
        } catch (NumberFormatException ex) {
            return false;
        }
    }
}
