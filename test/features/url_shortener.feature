Feature: URL Shortener
  As a user
  I want to shorten long URLs
  So that I can share them easily

  Background:
    Given the URL shortener service is running
    And I am on the home page

  Scenario Outline: Shorten valid URLs
    When I enter the URL "<url>"
    And I submit the form
    Then I should see a success message
    And I should see a shortened URL
    And the shortened URL should be valid
    And the shortned URL must redirect me to "<url>"

    Examples:
      | url                                                           |
      | https://github.com/ItsDobiel/URLShortener/blob/main/README.md |
      | https://cucumber.io/docs/bdd/                                 |

  Scenario Outline: URL normalization - trailing slashes
    When I enter the URL "<url_with_slash>"
    And I submit the form
    Then I should see a shortened URL
    When I enter the URL "<url_without_slash>"
    And I submit the form
    Then I should receive the same short code as before

    Examples:
      | url_with_slash                             | url_without_slash                         |
      | https://gorm.io/docs/                      | https://gorm.io/docs                      |
      | https://github.com/ItsDobiel/URLShortener/ | https://github.com/ItsDobiel/URLShortener |

  Scenario Outline: URL normalization - case insensitivity in protocol and domain
    When I enter the URL "<url_variant1>"
    And I submit the form
    Then I should see a shortened URL
    When I enter the URL "<url_variant2>"
    And I submit the form
    Then I should receive the same short code as before

    Examples:
      | url_variant1                              | url_variant2                              |
      | HTTPS://GiThUb.CoM/ItsDobiel/URLShortener | https://github.com/ItsDobiel/URLShortener |
      | httpS://docs.PODMAN.io/en/latest          | https://docs.podman.io/en/latest          |

  Scenario Outline: Duplicate URL returns same short code
    When I enter the URL "<url>"
    And I submit the form
    Then I should see a shortened URL
    When I enter the URL "<url>"
    And I submit the form
    Then I should receive the same short code as before

    Examples:
      | url                                                |
      | https://www.gnu.org/software/bash/manual/bash.html |
      | https://www.mozilla.org/en-US/about/manifesto      |

  Scenario Outline: Shorten URLs with query parameters, fragments and special characters
    When I enter the URL "<url>"
    And I submit the form
    Then I should see a success message
    And I should see a shortened URL

    Examples:
      | url                                                                |
      | https://example.com/search?q=Price%20of%20US%20%24&page=2#results  |
      | https://example.com/post?uuid=c3191902-bbef-4434-b043-e2cceedcc227 |

  Scenario Outline: Reject invalid URLs
    When I enter the URL "<invalid_url>"
    And I submit the form
    Then I should see an error message

    Examples:
      | invalid_url             |
      | not-a-valid-url         |
      | ftp://example.com/file  |
      | javascript:alert('xss') |
      | file:///etc/passwd      |
      | /../../../etc/passwd    |
      | /invalid@chars!         |
      | /short code with spaces |

  Scenario Outline: Accessing non-existent short code
    When I navigate to "<path>"
    Then I should see an error page
    And the error should indicate the short code was not found

    Examples:
      | path        |
      | /AuEcAowlXq |
      | /XXXXXXX    |

  Scenario Outline: Generated short code meets requirements
    When I enter the URL "<url>"
    And I submit the form
    Then I should see a shortened URL
    And the short code should be alphanumeric with allowed characters
    And the short code length should match the configured length

    Examples:
      | url                                             |
      | https://github.com/ItsDobiel/URLShortener       |
      | https://docs.fedoraproject.org/en-US/containers |
