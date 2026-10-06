---
title: 应用场景
layout: none
template: true
---

{{$solutions := pages "/solutions/*"}}
<section class="Solutions-header">
  <div class="Container">
    {{breadcrumbs .}}
    <h1 class="Solutions-useCase__title">应用场景</h1>
  </div>
</section>
<section class="Solutions-useCases">
  <div class="Container">
    <ul class="MarketingCardList">
      {{- range newest $solutions}}
        {{- if eq .series "Use Cases"}}
        <li class="MarketingCard">
          <a href="{{.URL}}">
            <div class="MarketingCard-section">
              {{- $icon := .icon}}
              {{- $iconDark := .iconDark}}
              {{- if $icon}}
              <img
                class="LightMode-img"
                loading="lazy"
                alt="{{$icon.alt}}"
                src="{{path.Dir .URL}}/{{$icon.file}}"
              />
              {{- end}}
              {{- if $iconDark}}
              <img
                class="DarkMode-img"
                loading="lazy"
                alt="{{$iconDark.alt}}"
                src="{{path.Dir .URL}}/{{$iconDark.file}}"
              />
              {{- end}}
            </div>
            <div class="MarketingCard-section MarketingCard-section__spacer">
              <h3 class="MarketingCard-title">{{or .linkTitle .title}}</h3>
              <p class="MarketingCard-body">
                {{.description}}
              </p>
            </div>
            <div class="MarketingCard-section__bottom" aria-describedby="usecase-description">
              <p class="MarketingCard-action">
                了解更多
              </p>
            </div>
          </a>
        </li>
        {{- end}}
      {{- end}}
    </ul>
  </div>
  <div class="screen-reader-only" id="usecase-description" hidden>
          在新窗口中打开。
    </div>
</section>
