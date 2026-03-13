package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.User;
import com.example.curs3projectback.model.enums.UserRole;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByTenant_CodeIgnoreCaseAndUsernameIgnoreCase(String tenantCode, String username);
    Optional<User> findByTenantIsNullAndUsernameIgnoreCase(String username);
    List<User> findAllByTenant_IdOrderByFullNameAsc(Long tenantId);
    long countByTenant_Id(Long tenantId);
    long countByRole(UserRole role);
    boolean existsByTenant_IdAndUsernameIgnoreCase(Long tenantId, String username);
}

