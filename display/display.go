package display

import (
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"training.go/Grypt/crypt"
)

type Display struct {
	app fyne.App

	window        fyne.Window
	windowEncrypt fyne.Window
	windowDecrypt fyne.Window

	buttonEncrypt              *widget.Button
	buttonEncryptValider       *widget.Button
	buttonCreateEncryptValider *widget.Button
	buttonDecrypt              *widget.Button
	buttonDecryptValider       *widget.Button
	buttonCopy                 *widget.Button
	buttonSecretCopy           *widget.Button
	buttonBytesCopy            *widget.Button
	buttonPwdEncryptCopy       *widget.Button

	radioEncrypt *widget.RadioGroup

	inputSecret          *widget.Entry
	inputByte            *widget.Entry
	inputPassword        *widget.Entry
	inputPasswordEncrypt *widget.Entry

	password   *canvas.Text
	err        *canvas.Text
	secret     *canvas.Text
	bytes      *canvas.Text
	encryptPwd *canvas.Text

	pwdBox *fyne.Container

	vBoxExistEncrypt  *fyne.Container
	vBoxCreateEncrypt *fyne.Container
}

type Displayer interface {
	createDecryptWindow()
	createEncryptWindow()
	createButton()
	createRadio()
	createInput()
	createText()

	decryptWindow()
	encryptWindow()

	encryptAction()
	decryptAction()
	radioEncryptAction(s string)
	validerAction()
	validerEncryptAction()
	validerCreateEncryptAction()

	cleanInput()
	printNewPwd()
}

var _ Displayer = (*Display)(nil)

func Run(a fyne.App) {
	w := a.NewWindow("Grypt")
	w.Resize(fyne.NewSize(200, 200))
	w.SetMaster()

	d := Display{
		app:    a,
		window: w,
	}

	d.createButton()
	d.createInput()
	d.createRadio()
	d.createText()

	title := canvas.NewText("Qu'est ce que vous voulez faire ?", color.White)

	hBox := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(), d.buttonEncrypt, d.buttonDecrypt, layout.NewSpacer())

	vBox := container.New(layout.NewVBoxLayout(), title, hBox)

	w.SetContent(vBox)
	w.ShowAndRun()
}

func (d *Display) createDecryptWindow() {
	wD := d.app.NewWindow("Dechiffré")
	wD.Resize(fyne.NewSize(1000, 700))
	wD.Hide()
	d.windowDecrypt = wD
}

func (d *Display) createEncryptWindow() {
	wE := d.app.NewWindow("Chiffré")
	wE.Resize(fyne.NewSize(1000, 700))
	wE.Hide()
	d.windowEncrypt = wE
}

func (d *Display) encryptWindow() {

	gridBox := container.New(layout.NewGridLayoutWithColumns(5),
		d.inputSecret, d.inputByte, d.inputPassword, layout.NewSpacer(), d.buttonEncryptValider)

	vBoxExist := container.New(layout.NewVBoxLayout(), gridBox,
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), d.password, d.buttonCopy, layout.NewSpacer()))

	vBoxCreate := container.New(layout.NewVBoxLayout(),
		d.err,
		container.New(layout.NewGridLayoutWithColumns(4), layout.NewSpacer(), d.inputPasswordEncrypt, d.buttonCreateEncryptValider, layout.NewSpacer()),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), widget.NewLabel("Le mot de passe :"), d.encryptPwd, d.buttonPwdEncryptCopy, layout.NewSpacer()),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), widget.NewLabel("Le secret :"), d.secret, d.buttonSecretCopy, layout.NewSpacer()),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), widget.NewLabel("Les bytes :"), d.bytes, d.buttonBytesCopy, layout.NewSpacer()),
	)

	vBox := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), d.radioEncrypt, layout.NewSpacer()),
		vBoxExist,
		vBoxCreate,
	)
	vBoxCreate.Hide()
	vBoxExist.Hide()

	d.vBoxCreateEncrypt = vBoxCreate
	d.vBoxExistEncrypt = vBoxExist

	d.windowEncrypt.SetContent(vBox)
}

func (d *Display) decryptWindow() {

	gridBox := container.New(layout.NewGridLayoutWithColumns(5),
		d.inputSecret, d.inputByte, d.inputPassword, layout.NewSpacer(), d.buttonDecryptValider)
	pwdBox := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), d.password, d.buttonCopy, layout.NewSpacer())

	vBox := container.New(layout.NewVBoxLayout(), d.err, gridBox, pwdBox)

	d.pwdBox = pwdBox

	d.windowDecrypt.SetContent(vBox)
}

func (d *Display) createButton() {
	// Manage page
	bD := widget.NewButton("Dechiffré", d.decryptAction)
	bE := widget.NewButton("Chiffré", d.encryptAction)

	// Decrypt page
	bDV := widget.NewButton("Valider", d.validerAction)

	// Encrypt page
	bEV := widget.NewButton("Valider", d.validerEncryptAction)
	bCEV := widget.NewButton("Valider", d.validerCreateEncryptAction)

	bC := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		d.window.Clipboard().SetContent(d.password.Text)
	})
	bCSecret := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		d.window.Clipboard().SetContent(d.secret.Text)
	})
	bCPwdE := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		d.window.Clipboard().SetContent(d.encryptPwd.Text)
	})
	bCByte := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		d.window.Clipboard().SetContent(d.bytes.Text)
	})

	bCSecret.Hide()
	bCPwdE.Hide()
	bCByte.Hide()
	bC.Hide()

	d.buttonEncryptValider = bEV
	d.buttonCreateEncryptValider = bCEV
	d.buttonCopy = bC
	d.buttonDecryptValider = bDV
	d.buttonEncrypt = bE
	d.buttonDecrypt = bD
	d.buttonPwdEncryptCopy = bCPwdE
	d.buttonSecretCopy = bCSecret
	d.buttonBytesCopy = bCByte
}

func (d *Display) createInput() {
	inputS := widget.NewEntry()
	inputS.SetPlaceHolder("Entrez le secret...")

	inputB := widget.NewEntry()
	inputB.SetPlaceHolder("Entrez les bytes 00,00,00...")

	inputPwd := widget.NewEntry()
	inputPwd.SetPlaceHolder("Entre le mot de passe")

	inputPwdE := widget.NewEntry()
	inputPwdE.SetPlaceHolder("Entre le mot de passe")

	d.inputPassword = inputPwd
	d.inputSecret = inputS
	d.inputByte = inputB
	d.inputPasswordEncrypt = inputPwdE
}

func (d *Display) createRadio() {
	rEncrypt := widget.NewRadioGroup([]string{
		" - Utiliser un secret existant",
		" - Créer un nouveau secret"},
		d.radioEncryptAction)

	d.radioEncrypt = rEncrypt
}

func (d *Display) createText() {
	pwd := canvas.NewText("", color.White)
	err := canvas.NewText("", color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 1,
	})

	secret := canvas.NewText("", color.White)
	bytes := canvas.NewText("", color.White)
	encryptPwd := canvas.NewText("", color.White)

	d.secret = secret
	d.bytes = bytes
	d.encryptPwd = encryptPwd
	d.password = pwd
	d.err = err
}

func (d *Display) decryptAction() {
	d.createDecryptWindow()
	d.windowDecrypt.Show()
	if d.windowEncrypt != nil {
		d.windowEncrypt.Close()
	}
	d.cleanInput()
	d.decryptWindow()
}

func (d *Display) encryptAction() {
	d.createEncryptWindow()
	d.windowEncrypt.Show()
	if d.windowDecrypt != nil {
		d.windowDecrypt.Close()
	}
	d.cleanInput()
	d.encryptWindow()
}

func (d *Display) radioEncryptAction(s string) {
	d.cleanInput()
	if strings.Contains(s, "Utiliser un secret") {
		d.vBoxExistEncrypt.Show()
		d.vBoxCreateEncrypt.Hide()
	} else {
		d.vBoxCreateEncrypt.Show()
		d.vBoxExistEncrypt.Hide()
	}
}

func (d *Display) validerAction() {
	if len(d.inputByte.Text) == 0 && len(d.inputSecret.Text) == 0 && len(d.inputPassword.Text) == 0 {
		d.err.Text = "S'il vous plait rempliser les 3 champs"
		d.err.Show()
	} else {
		d.err.Hide()
	}
	b := strings.Split(d.inputByte.Text, ",")
	bytes, err := stringSliceToByteSlice(b)
	if err != nil {
		log.Print(err)
		return
	}

	p, e := crypt.Decrypt(d.inputPassword.Text, d.inputSecret.Text,
		bytes)
	if e != nil {
		log.Println(e)
	}
	d.password.Text = p
	d.buttonCopy.Show()
	d.password.Show()
}

func (d *Display) validerEncryptAction() {
	if len(d.inputByte.Text) == 0 && len(d.inputSecret.Text) == 0 && len(d.inputPassword.Text) == 0 {
		d.err.Text = "S'il vous plait rempliser les 3 champs"
		d.err.Show()
	} else {
		d.err.Hide()
	}
	b := strings.Split(d.inputByte.Text, ",")
	bytes, err := stringSliceToByteSlice(b)
	if err != nil {
		log.Print(err)
		return
	}

	p, err := crypt.Encrypt(
		d.inputPassword.Text,
		d.inputSecret.Text,
		bytes)

	if err != nil {
		log.Println(err)
	}
	d.password.Text = p
	d.buttonCopy.Show()
	d.password.Show()
}

func (d *Display) cleanInput() {
	d.inputByte.SetText("")
	d.inputPassword.SetText("")
	d.inputSecret.SetText("")
	d.password.Text = ""
	d.radioEncrypt.Selected = ""
	d.password.Hide()
	d.buttonCopy.Hide()

	d.buttonBytesCopy.Hide()
	d.buttonSecretCopy.Hide()
	d.buttonPwdEncryptCopy.Hide()

	d.secret.Text = ""
	d.bytes.Text = ""
	d.encryptPwd.Text = ""
}

func (d *Display) validerCreateEncryptAction() {
	if len(d.inputPasswordEncrypt.Text) <= 0 {
		d.err.Text = "S'il vous plaît remplissez le champs mot de passe."
		d.err.Show()
		return
	}
	d.printNewPwd()
}

func (d *Display) printNewPwd() {
	secret, byt, encryptPwd, err := generateNewPwd(d.inputPasswordEncrypt.Text)
	if err != nil {
		d.err.Text = err.Error()
		return
	}
	d.secret.Text = secret
	d.bytes.Text = string(byt)
	d.encryptPwd.Text = encryptPwd
	d.buttonBytesCopy.Show()
	d.buttonSecretCopy.Show()
	d.buttonPwdEncryptCopy.Show()
}
