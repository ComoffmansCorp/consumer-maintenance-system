package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.dictionary.AddressResponse;
import com.example.curs3projectback.dto.dictionary.OrganizationResponse;
import com.example.curs3projectback.mapper.DictionaryMapper;
import com.example.curs3projectback.repository.AddressRepository;
import com.example.curs3projectback.repository.OrganizationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class DictionaryService {

    private final AddressRepository addressRepository;
    private final OrganizationRepository organizationRepository;
    private final DictionaryMapper dictionaryMapper;
    private final CurrentTenantService currentTenantService;

    public List<AddressResponse> findAddresses(String query, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Searching addresses tenant={} query={}", tenant.getCode(), query);
        return dictionaryMapper.toAddressList(addressRepository.findAllByTenantAndStreetContainingIgnoreCaseOrderByStreet(tenant, query));
    }

    public List<OrganizationResponse> findOrganizations(String query, Authentication authentication) {
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Searching organizations tenant={} query={}", tenant.getCode(), query);
        return dictionaryMapper.toOrganizationList(organizationRepository.findAllByTenantAndNameContainingIgnoreCaseOrderByName(tenant, query));
    }
}
