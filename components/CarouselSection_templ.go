// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/jmarren/deepfried/services"

func CarouselSectionBody(carousel []*services.PlayableElt) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card-carousel\" data-click-index=\"0\" hx-on:htmx-after-settle=\"triggerCardAnimations(this)\"><div class=\"cards-container\" data-current-transform=\"0\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, card := range carousel {
			templ_7745c5c3_Err = CarouselCardBody(card).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button class=\"carousel-chevron-forward\" hx-on:click=\"MoveCarouselRight(this)\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ChevronRightIcon().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var carouselHeadOnceHandle = templ.NewOnceHandle()

func CarouselSectionHead() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var3 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script id=\"carousel-script\">\n\tfunction MoveCarouselRight(elt) {\n\t\tconst carouselElt = htmx.closest(elt, \".card-carousel\")\n\t\tconst carouselWidth = carouselElt.getBoundingClientRect().width \n\t\tconst cardsContainerElts = Array.from(htmx.findAll(carouselElt, \".cards-container\"))\n\t\tlet cardsContainerElt\n\t\tif (cardsContainerElts.length > 1) {\n\t\t\tcardsContainerElts[0].remove()\n\t\t\tcardsContainerElt = cardsContainerElts[1]\n\t\t} else {\n\t\t\tcardsContainerElt = cardsContainerElts[0]\n\t\t}\n\t\tconst cardElt = htmx.find(cardsContainerElt, \".track-card-container\")\n\t\tconst cardWidth = cardElt.getBoundingClientRect().width\n\t\tconst cardsArr = Array.from(htmx.findAll(cardsContainerElt, \".track-card-container\"))\n\t\tconst cardsContainerWidth = cardsContainerElt.getBoundingClientRect().width\n\t\tconst slideDist = -1 * carouselWidth\n\t\tlet currentTransform = parseInt(cardsContainerElt.dataset.currentTransform)\n\t\tconst totalWidth = cardsArr.length * cardWidth\n\t\tconst paddedWidth = totalWidth - carouselWidth\n\n\t\tif (Math.abs(currentTransform) > paddedWidth) { \n\t\t\tconst newCards = cardsContainerElt.cloneNode(true)\n\t\t\tcardsContainerElt.style.position = 'absolute';\n\t\t\tnewCards.dataset.currentTransform = 0\n\t\t\tcarouselElt.appendChild(newCards)\n\t\t\thtmx.process(newCards)\n\t\t\t\n\t\t\tconst frames = [\n\t\t\t\t{ transform: `translateX(${carouselWidth}px)`},\n\t\t\t\t{ transform: `translateX(${0}px)`}\n\t\t\t];\n\t\t\tconst animation = new KeyframeEffect(\n\t\t\t\tnewCards,\n\t\t\t\tframes, \n\t\t\t\t{ duration: 500, fill: 'forwards', easing: 'cubic-bezier(0.15, 0.83, 0.66, 1)'}\n\t\t\t\t)\n\t\t\tconst player = new Animation(animation, document.timeline)\n\t\t\tplayer.play()\n\t\t} \n\n\t\tconst frames = [\n\t\t\t{ transform: `translateX(${currentTransform}px)`},\n\t\t\t{ transform: `translateX(${currentTransform + slideDist + cardWidth}px)`}\n\t\t];\n\t\tconst animation = new KeyframeEffect(\n\t\t\tcardsContainerElt,\n\t\t\tframes, \n\t\t\t{ duration: 500, fill: 'forwards', easing: 'cubic-bezier(0.15, 0.83, 0.66, 1)'}\n\t\t\t)\n\t\tconst player = new Animation(animation, document.timeline)\n\t\tplayer.play()\n\t\tcardsContainerElt.dataset.currentTransform = currentTransform + slideDist\n\t}\n\n\t</script> <style id=\"carouselStyles\">\n  .card-carousel {\n    min-width: 500px;\n    flex-grow: 1;\n    padding: 10px 10px 15px 10px;\n    overflow-x: hidden;\n    position: relative;\n    margin-bottom: 10px;\n    margin-top: 5px;\n/*    background-color: #1d1d1e;*/\n/*   background-color: var(--dark-gray); */\n/*    background-color: #282727; */\n     border-radius: 10px 0 0 10px;\n/*     background-color: #1b1f20; */\n/*    background-color: var(--dark-gray); */\n  }\n\n   .cards-container {\n    flex-grow: 1;\n    display: flex;\n    flex-direction: row;\n    gap: 10px;\n    position: inline-block;\n   }\n\n\n   @keyframes slide-left {\n\t0% {\n\t\ttransform: 0;\n\t}\n\t100% {\n\t\ttransform: translateX(calc(-100vw + 100px));\n\t}\n   }\n\n\n  .card-carousel .track-image {\n    width: var(--medium-image-size);\n    height: var(--medium-image-size);\n    min-width: var(--medium-image-size);\n    min-height: var(--medium-image-size);\n    border-radius: 10px;\n  }\n\n  .carousel-chevron-forward {\n\twidth: 20px;\n\theight: 20px;\n\tborder-radius: 100%;\n\tposition: absolute;\n\tright: 10px;\n\ttop: calc(50% - 10px);\n\tbackground-color: white;\n\tdisplay: flex;\n\tjustify-content: center;\n\talign-items: center;\n\tz-index: 3;\n\tborder: 1px solid var(--border-color-1);\n  }\n\n  .carousel-chevron-forward:hover {\n\tbackground-color: var(--border-color-1);\n  }\n  .card-carousel .track-image>img {\n    border-radius: 5px;\n  }\n</style>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = carouselHeadOnceHandle.Once().Render(templ.WithChildren(ctx, templ_7745c5c3_Var3), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
