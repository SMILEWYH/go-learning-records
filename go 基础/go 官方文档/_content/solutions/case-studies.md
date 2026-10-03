---
title: 案例研究
layout: none
template: true
---

{{$solutions := pages "/solutions/*"}}
<section class="Solutions-headline">
  <div class="GoCarousel" id="SolutionsHeroCarousel-carousel">
    <div class="GoCarousel-controlsContainer">
      <div class="GoCarousel-wrapper SolutionsHeroCarousel-wrapper">
      {{ breadcrumbs . }}
        <ul class="js-solutionsHeroCarouselSlides SolutionsHeroCarousel-slides">
          {{- $n := 0}}
          {{- range newest $solutions}}
            {{- if eq .series "Case Studies"}}
              {{- $n = add $n 1}}
              {{- if le $n 3}}
              <li class="SolutionsHeroCarousel-slide">
                <div class="Solutions-headlineImg">
                  <img
                    src="/images/{{.carouselImgSrc}}"
                    alt="{{(or .linkTitle .title)}}"
                  />
                </div>
                <div class="Solutions-headlineText">
                  <p class="Solutions-headlineNotification">最近更新</p>
                  <h2>
                    {{(or .linkTitle .title)}}
                  </h2>
                  <p class="Solutions-headlineBody">
                    {{with .quote}}{{.}}{{end}}
                    <a href="{{.URL}}" aria-describedby="casestudy-description"
                      >了解更多
                      <i class="material-icons Solutions-forwardArrowIcon"
                         aria-hidden="true"><svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" focusable="false"><path d="M4 12h16M13 5l7 7-7 7"/></svg></i>
                    </a>
                  </p>
                </div>
              </li>
              {{- end}}
            {{- end}}
          {{- end}}
        </ul>
      </div>
      <div class="screen-reader-only" id="casestudy-description" hidden>
          在新窗口中打开。
    </div>
      <button
        class="js-solutionsHeroCarouselPrev GoCarousel-controlPrev GoCarousel-controlPrev-solutionsHero"
        hidden
      >
        <i class="GoCarousel-icon material-icons" aria-hidden="true"><svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" focusable="false"><path d="M15 6l-6 6 6 6"/></svg></i>
      </button>
      <button
        class="js-solutionsHeroCarouselNext GoCarousel-controlNext GoCarousel-controlNext-solutionsHero"
      >
        <i class="GoCarousel-icon material-icons" aria-hidden="true"><svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" focusable="false"><path d="M9 6l6 6-6 6"/></svg></i>
      </button>
    </div>
  </div>
</section>

<section class="Solutions-useCases">
  <div class="Container">
    <ul class="MarketingCardList">
      {{- range $solutions}}
        {{- if eq .series "Case Studies"}}
          <li class="MarketingCard">
            {{- if .link}}
            <a
            href="{{.link}}"
            target="_blank"
            rel="noopener"
            >
            {{- else}}
            <a href="{{.URL}}">
            {{- end}}
              <div class="MarketingCard-section">
                <img
                  height="36"
                  class="DarkMode-img"
                  loading="lazy"
                  alt="{{.company}}"
                  src="/images/logos/{{.logoSrcDark}}"
                />
                <img
                  height="36"
                  class="LightMode-img"
                  loading="lazy"
                  alt="{{.company}}"
                  src="/images/logos/{{.logoSrc}}"
                />
              </div>
              <div class="MarketingCard-section MarketingCard-section__spacer">
                <h2 class="MarketingCard-title">{{or .linkTitle .title}}</h2>
                <p class="MarketingCard-body">
                {{- if .link}}
                  {{.description}}
                {{- else}}
                  {{with .quote}}{{.}}{{end}}
                {{- end}}
                </p>
              </div>
              <div class="MarketingCard-section__bottom" aria-describedby="casestudy-description">
                <p class="MarketingCard-action">查看案例</p>
              </div>
            </a>
          </li>
        {{- end}}
      {{- end}}
    <ul>
  </div>
</section>
